"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var mobx_1 = require("mobx");
function printDebugValue(v) {
    return mobx_1.getDependencyTree(v);
}
exports.printDebugValue = printDebugValue;
