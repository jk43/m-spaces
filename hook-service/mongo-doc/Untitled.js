/*
 Navicat Premium Data Transfer

 Source Server         : k3s mongodb - DEV
 Source Server Type    : MongoDB
 Source Server Version : 60002 (6.0.2)
 Source Host           : 192.168.1.55:27017
 Source Schema         : dev-hook

 Target Server Type    : MongoDB
 Target Server Version : 60002 (6.0.2)
 File Encoding         : 65001

 Date: 30/01/2024 15:29:05
*/


// ----------------------------
// Collection structure for hooks
// ----------------------------
db.getCollection("hooks").drop();
db.createCollection("hooks");

// ----------------------------
// Documents of hooks
// ----------------------------
db.getCollection("hooks").insert([ {
    _id: ObjectId("65b95aad4becb142e2024874"),
    host: "localhost",
    method: "put",
    path: "/user/user",
    hookType: "rest",
    hookPosition: "pre",
    hookAddress: "hook-service",
    hookFunction: "localhost/pre/put/user",
    description: ""
}, {
    _id: ObjectId("65b95b244becb142e2024875"),
    host: "localhost",
    method: "put",
    path: "/user/user",
    hookType: "rest",
    hookPosition: "post",
    hookAddress: "hook-service",
    hookFunction: "localhost/post/put/user"
} ]);
