"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var mobx_1 = require("mobx");
var react_1 = require("react");
if (!react_1.useState) {
    throw new Error("mobx-react-lite requires React with Hooks support");
}
if (!mobx_1.spy) {
    throw new Error("mobx-react-lite requires mobx at least version 4 to be available");
}
