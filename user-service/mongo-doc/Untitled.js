/*
 Navicat Premium Data Transfer

 Source Server         : k3s mongodb - DEV
 Source Server Type    : MongoDB
 Source Server Version : 60002 (6.0.2)
 Source Host           : 192.168.1.55:27017
 Source Schema         : dev-user

 Target Server Type    : MongoDB
 Target Server Version : 60002 (6.0.2)
 File Encoding         : 65001

 Date: 30/01/2024 15:29:59
*/


// ----------------------------
// Collection structure for users
// ----------------------------
db.getCollection("users").drop();
db.createCollection("users");

// ----------------------------
// Documents of users
// ----------------------------
db.getCollection("users").insert([ {
    _id: ObjectId("650b4efdf04d408b0538309f"),
    email: "jk@jktech.net",
    organizations: [
        {
            _id: ObjectId("64540b79cb7c89d64e01ff45"),
            firstName: "hyojun",
            lastName: "kim",
            role: "member",
            status: "active",
            metadata: [
                {
                    key: "_createdBy",
                    value: "alex kim"
                },
                {
                    key: "_level",
                    value: 1
                },
                {
                    key: "newsletter",
                    value: "dddd1111"
                },
                {
                    key: "phone",
                    value: "gogo hello"
                },
                {
                    key: "zipcode",
                    value: "zip"
                },
                {
                    key: "notification",
                    value: "ccc yyy"
                },
                {
                    key: "gender",
                    value: "gender1111"
                },
                {
                    key: "groups",
                    value: "xyz kkk"
                }
            ],
            created_at: ISODate("2023-09-20T19:58:53.896Z")
        },
        {
            _id: ObjectId("64cc081d96abe03a3b060ea9"),
            firstName: "alex",
            lastName: "kim",
            role: "member",
            status: "active",
            metadata: [
                {
                    key: "deleteme1",
                    value: "alex kim"
                },
                {
                    key: "deleteme2",
                    value: "alex kim"
                },
                {
                    key: "deleteme3",
                    value: "alex kim"
                },
                {
                    key: "deleteme4",
                    value: "alex kim"
                }
            ],
            created_at: ISODate("2023-09-20T19:59:21.181Z")
        }
    ]
}, {
    _id: ObjectId("65424f4a1687ab5d3810c74a"),
    email: "mama@mama.com",
    organizations: [
        {
            _id: ObjectId("64540b79cb7c89d64e01ff45"),
            firstName: "gogo",
            lastName: "mama",
            role: "member",
            status: "active",
            metadata: [
                {
                    key: "deleteme1",
                    value: "alex kim"
                },
                {
                    key: "deleteme2",
                    value: "alex kim"
                },
                {
                    key: "deleteme3",
                    value: "alex kim"
                },
                {
                    key: "deleteme4",
                    value: "alex kim"
                }
            ],
            created_at: ISODate("2023-11-01T13:14:50.622Z")
        }
    ]
}, {
    _id: ObjectId("65761f92d43e101ec019d778"),
    email: "jame@jktech.net",
    organizations: [
        {
            _id: ObjectId("64540b79cb7c89d64e01ff45"),
            firstName: "james",
            lastName: "kim",
            role: "member",
            status: "active",
            metadata: [
                {
                    key: "deleteme1",
                    value: "alex kim"
                },
                {
                    key: "deleteme2",
                    value: "alex kim"
                },
                {
                    key: "deleteme3",
                    value: "alex kim"
                },
                {
                    key: "deleteme4",
                    value: "alex kim"
                }
            ],
            created_at: ISODate("2023-12-10T20:29:06.306Z")
        }
    ]
} ]);
