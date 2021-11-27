"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var globalIsUsingStaticRendering = false;
function useStaticRendering(enable) {
    globalIsUsingStaticRendering = enable;
}
exports.useStaticRendering = useStaticRendering;
function isUsingStaticRendering() {
    return globalIsUsingStaticRendering;
}
exports.isUsingStaticRendering = isUsingStaticRendering;
