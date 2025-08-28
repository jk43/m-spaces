package service

import (
	"context"
	"os"

	"github.com/moly-space/molylibs/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func DeleteFile(userID, orgID, uuID string) (string, error) {
	conn, err := grpc.Dial(os.Getenv("FILE_SERVICE_ADDR")+":5000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		return "", err
	}
	c := pb.NewFileServiceClient(conn)
	// sets := make([]string, len(metadataSets))
	// for i, p := range metadataSets {
	// 	sets[i] = p.Act
	// }
	in := pb.DeleteFileRequest{
		UserID: userID,
		OrgID:  orgID,
		UuID:   uuID,
	}
	res, err := c.DeleteFile(context.Background(), &in)
	if err != nil {
		return "", err
	}
	return res.Id, nil
}
