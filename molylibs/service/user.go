package service

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/moly-space/molylibs/pb"
	"github.com/moly-space/molylibs/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserStatus string
type RegisterMethod string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusDeleted  UserStatus = "deleted"
	UserStatusWaiting  UserStatus = "waiting"
)

const (
	RegisterMethodAdmin         RegisterMethod = "admin"
	RegisterMethodSelf          RegisterMethod = "self"
	RegisterMethodOAuthGoogle   RegisterMethod = "oauth_google"
	RegisterMethodOAuthApple    RegisterMethod = "oauth_apple"
	RegisterMethodOAuthFacebook RegisterMethod = "oauth_facebook"
)

type Metadata = utils.MapStringAny
type StoreData struct {
	Key   string             `bson:"key" json:"key"`
	Label string             `bson:"label" json:"label"`
	Value utils.MapStringAny `bson:"value" json:"value"`
}
type User struct {
	ID             primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	OrganizationID primitive.ObjectID     `bson:"organizationId" json:"org_id"`
	Email          string                 `bson:"email" json:"email" validate:"required,email"`
	FirstName      string                 `bson:"firstName" json:"firstName,omitempty" validate:"required"`
	LastName       string                 `bson:"lastName" json:"lastName,omitempty" validate:"required"`
	ProfileImage   string                 `bson:"profileImage" json:"profileImage,omitempty"`
	Role           utils.UserRole         `bson:"role" json:"role"`
	Status         UserStatus             `bson:"status" json:"status"`
	Metadata       map[string]Metadata    `bson:"metadata" json:"metadata"`
	CreatedAt      time.Time              `bson:"createdAt" json:"createdAt"`
	Store          map[string][]StoreData `bson:"store" json:"store"`
	Verified       bool                   `bson:"verified" json:"verified"`
	//SelfRegistered bool                   `bson:"selfRegistered" json:"selfRegistered"`
	RegisterMethod RegisterMethod `bson:"registerMethod" json:"registerMethod"`
}

// This method parses metadata from the User document.
// The 'sets' parameter represents the objects within the metadata that we want to parse.
// For example, to parse metadata for the 'user' object, we would call the function with '[]string{service.FormNameForUserSettings}' as the value of the 'sets' parameter.
// The 'filter' parameter represents the pattern of keys in the selected object that we want to exclude for security reasons.
// For example, to exclude elements with keys like '__sensitiveInfo', we would call the function with '__' as the filter value.
func (u *User) GetMetadata(forms []string, shareable bool) (map[string]Metadata, error) {
	orgSettings, err := GetOrg(nil, u.OrganizationID.Hex())
	if err != nil {
		return nil, err
	}
	metadata := make(map[string]Metadata)
	// Which metadata to parse
	for _, f := range forms {
		metadata[f] = make(Metadata)
		if shareable {
			for _, i := range orgSettings.Forms[f] {
				if shareable && !i.Shareable {
					continue
				}
				metadata[f][i.Key] = u.Metadata[f][i.Key]
			}
		} else {
			for k, m := range u.Metadata[f] {
				metadata[f][k] = m
			}
		}
	}
	return metadata, nil
}

func GetUser(id, email, orgID string, editableMetadataOnly, shareableMetadataOnly bool, forms []string) (*User, error) {
	conn, err := grpc.Dial(os.Getenv("USER_SERVICE_ADDR")+":5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		return nil, err
	}
	c := pb.NewUserServiceClient(conn)
	// sets := make([]string, len(metadataSets))
	// for i, p := range metadataSets {
	// 	sets[i] = p.Act
	// }
	in := pb.UserRequest{
		UserID:                id,
		Email:                 email,
		EditableMetadataOnly:  editableMetadataOnly,
		ShareableMetadataOnly: shareableMetadataOnly,
		MetadataSets:          forms,
	}
	res, err := c.GetUser(context.Background(), &in)
	if err != nil {
		return nil, err
	}
	userId, err := primitive.ObjectIDFromHex(res.ID)
	userOrgID, err := primitive.ObjectIDFromHex(res.OrgID)
	user := &User{
		ID:             userId,
		OrganizationID: userOrgID,
		Email:          res.Email,
		FirstName:      res.FirstName,
		LastName:       res.LastName,
		ProfileImage:   res.ProfileImage,
		Role:           utils.UserRole(res.Role),
		Status:         UserStatus(res.Status),
		Metadata:       make(map[string]Metadata),
		RegisterMethod: RegisterMethod(res.RegisterMethod),
	}
	err = json.Unmarshal(res.Metadata, &user.Metadata)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func ProcessOAuthUser(user *User) (*pb.UserResponse, error) {
	conn, err := grpc.Dial(os.Getenv("USER_SERVICE_ADDR")+":5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		return nil, err
	}
	c := pb.NewUserServiceClient(conn)
	in := pb.OAuthUserRequest{
		OrgID:          user.OrganizationID.Hex(),
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		ProfileImage:   user.ProfileImage,
		RegisterMethod: string(user.RegisterMethod),
	}
	res, err := c.ProcessOAuthUser(context.Background(), &in)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func GetUserWithRole(id, email, orgID, host string, service utils.Service, unit utils.Ctx, role utils.UserRole, editableMetadataOnly, shareableMetadataOnly bool) (*User, error) {
	perm, err := utils.NewCasbin(host, utils.ServiceOrg, utils.CtxOrgForm, "").ParsePermission(role, utils.GeneralPermissionParser)
	if err != nil {
		return nil, err
	}
	return GetUser(id, email, orgID, editableMetadataOnly, shareableMetadataOnly, perm)
}
