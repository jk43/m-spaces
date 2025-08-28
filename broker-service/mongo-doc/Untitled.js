/*
 Navicat Premium Data Transfer

 Source Server         : Docker moly mongo
 Source Server Type    : MongoDB
 Source Server Version : 70008 (7.0.8)
 Source Host           : localhost:27017
 Source Schema         : broker

 Target Server Type    : MongoDB
 Target Server Version : 70008 (7.0.8)
 File Encoding         : 65001

 Date: 28/07/2024 18:40:01
*/


// ----------------------------
// Collection structure for http_rules
// ----------------------------
db.getCollection("http_rules").drop();
db.createCollection("http_rules");

// ----------------------------
// Documents of http_rules
// ----------------------------
db.getCollection("http_rules").insert([ {
    _id: ObjectId("637660fd94a83f0adb0b09c3"),
    url: "http://user-service/save",
    method: "POST",
    path: "/user/save"
}, {
    _id: ObjectId("6376611394a83f0adb0b09c4"),
    url: "http://auth-service/login",
    method: "POST",
    path: "/user/login"
}, {
    _id: ObjectId("63af8a6d1a5cc6253d0728a2"),
    url: "http://auth-service/logout",
    method: "POST",
    path: "/user/logout"
}, {
    _id: ObjectId("63af8b851a5cc6253d0728a3"),
    method: "GET",
    path: "/user/refreshtoken",
    url: "http://auth-service/refreshtoken"
}, {
    _id: ObjectId("63b040ce1a5cc6253d0728a4"),
    url: "http://user-service/user",
    method: "GET",
    path: "/user/user",
    description: "Pulling user information to store in web"
}, {
    _id: ObjectId("63e2bb6b3377118d6109df22"),
    url: "http://auth-service/verifyemail",
    method: "POST",
    path: "/user/verifyemail"
}, {
    _id: ObjectId("64ff77892889dfcfb604aeba"),
    url: "http://sandbox-service/test",
    method: "GET",
    path: "/sandbox/test"
}, {
    _id: ObjectId("650062022889dfcfb604aebc"),
    method: "*",
    path: "/sandbox/*",
    url: "http://sandbox-service/"
}, {
    _id: ObjectId("651c61ff778fa4bd5b014901"),
    url: "http://organization-service/settings",
    method: "GET",
    path: "/organization/settings"
}, {
    _id: ObjectId("654a3928eddcc2bcb70fc293"),
    url: "http://organization-service/items",
    method: "GET",
    path: "/organization/items"
}, {
    _id: ObjectId("658f251bed3564c9b3004d63"),
    url: "http://user-service/user/metadata",
    method: "PUT",
    path: "/user/settings"
}, {
    _id: ObjectId("65d0d8b64becb142e202487e"),
    url: "http://user-service/user/account",
    method: "PUT",
    path: "/user/account"
}, {
    _id: ObjectId("65d0d8c54becb142e202487f"),
    url: "http://user-service/user/password",
    method: "PUT",
    path: "/user/password",
    description: "Update password and server will send verification code"
}, {
    _id: ObjectId("65e08014e7122d564603f432"),
    url: "http://sandbox-service/sandbox",
    method: "GET",
    path: "/sandbox/sandbox"
}, {
    _id: ObjectId("65f1f61fadf2cf52220c9093"),
    url: "http://user-service/user/password",
    method: "POST",
    path: "/user/password",
    description: "Update password with verification code"
}, {
    _id: ObjectId("6603318bec327629410f0832"),
    url: "http://user-service/user/email",
    method: "POST",
    description: "Update email with verification code",
    path: "/user/email"
}, {
    _id: ObjectId("66042460ec327629410f0835"),
    url: "http://auth-service/forgotpassword",
    method: "POST",
    path: "/auth/forgotpassword"
}, {
    _id: ObjectId("66042465ec327629410f0836"),
    method: "PUT",
    path: "/auth/password",
    url: "http://auth-service/password"
}, {
    _id: ObjectId("6604d20eec327629410f0839"),
    url: "http://file-service/upload",
    method: "POST",
    path: "/files"
}, {
    _id: ObjectId("662c3d744d14cabf2e0bf323"),
    url: "http://organization-service/info",
    method: "PUT",
    path: "/organization/info"
}, {
    _id: ObjectId("662c3d794d14cabf2e0bf324"),
    method: "GET",
    path: "/organization/info",
    url: "http://organization-service/info"
}, {
    _id: ObjectId("662d2f894d14cabf2e0bf326"),
    url: "http://user-service/admin/users",
    method: "GET",
    path: "/admin/users"
}, {
    _id: ObjectId("662ff6914d14cabf2e0bf32d"),
    url: "http://user-service/admin/user",
    method: "PUT",
    path: "/admin/user"
}, {
    _id: ObjectId("662ff6a24d14cabf2e0bf32e"),
    url: "http://user-service/admin/user",
    method: "POST",
    path: "/admin/user"
}, {
    _id: ObjectId("663212854d14cabf2e0bf334"),
    method: "DELETE",
    path: "/admin/user",
    url: "http://user-service/admin/user"
}, {
    _id: ObjectId("663213024d14cabf2e0bf335"),
    method: "PUT",
    path: "/admin/organization/form-order",
    url: "http://organization-service/form-order"
}, {
    _id: ObjectId("663213114d14cabf2e0bf336"),
    method: "PUT",
    path: "/admin/organization/form",
    url: "http://organization-service/form"
}, {
    _id: ObjectId("663213194d14cabf2e0bf337"),
    method: "DELETE",
    path: "/admin/organization/form",
    url: "http://organization-service/form"
}, {
    _id: ObjectId("663213204d14cabf2e0bf338"),
    method: "POST",
    path: "/admin/organization/form",
    url: "http://organization-service/form"
}, {
    _id: ObjectId("663bc7944d14cabf2e0bf33e"),
    url: "http://organization-service/form",
    method: "GET",
    path: "/admin/organization/form/*"
}, {
    _id: ObjectId("664c8ae06829f565220d59a0"),
    method: "POST",
    path: "/user/setpassword",
    url: "http://auth-service/setpassword"
}, {
    _id: ObjectId("665459c76829f565220d59a2"),
    url: "http://user-service/user/store",
    method: "POST",
    path: "/user/store"
}, {
    _id: ObjectId("66545a136829f565220d59a3"),
    url: "http://user-service/user/store",
    method: "PUT",
    path: "/user/store"
}, {
    _id: ObjectId("66545a236829f565220d59a4"),
    url: "http://user-service/user/store",
    method: "DELETE",
    path: "/user/store"
}, {
    _id: ObjectId("66545a3f6829f565220d59a5"),
    url: "http://user-service/user/store",
    method: "GET",
    path: "/user/store/*"
}, {
    _id: ObjectId("66843d8fd81a728793025399"),
    url: "http://tree-service/tree",
    method: "GET",
    path: "/tree/*"
}, {
    _id: ObjectId("66843dabd81a72879302539a"),
    url: "http://tree-service/admin/tree",
    method: "GET",
    path: "/admin/tree"
}, {
    _id: ObjectId("66843db8d81a72879302539b"),
    url: "http://tree-service/admin/tree",
    method: "PUT",
    path: "/admin/tree"
}, {
    _id: ObjectId("66843dc1d81a72879302539c"),
    url: "http://tree-service/admin/tree",
    method: "POST",
    path: "/admin/tree"
}, {
    _id: ObjectId("66843dd8d81a72879302539d"),
    url: "http://tree-service/admin/tree",
    method: "DELETE",
    path: "/admin/tree"
}, {
    _id: ObjectId("6687054cd81a7287930253a1"),
    url: "http://tree-service/admin/trees",
    method: "GET",
    path: "/admin/trees"
}, {
    _id: ObjectId("668d8edfd81a7287930253a2"),
    method: "PUT",
    path: "/admin/tree/reorder",
    url: "http://tree-service/admin/reorder"
}, {
    _id: ObjectId("6696e355d81a7287930253a4"),
    url: "http://tree-service/admin/trees",
    method: "POST",
    path: "/admin/trees"
} ]);
