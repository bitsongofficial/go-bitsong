"use strict";
var __read = (this && this.__read) || function (o, n) {
    var m = typeof Symbol === "function" && o[Symbol.iterator];
    if (!m) return o;
    var i = m.call(o), r, ar = [], e;
    try {
        while ((n === void 0 || n-- > 0) && !(r = i.next()).done) ar.push(r.value);
    }
    catch (error) { e = { error: error }; }
    finally {
        try {
            if (r && !r.done && (m = i["return"])) m.call(i);
        }
        finally { if (e) throw e.error; }
    }
    return ar;
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
var mobx_1 = require("mobx");
var react_1 = __importDefault(require("react"));
var utils_1 = require("./utils");
function useAsObservableSourceInternal(current, usedByLocalStore) {
    var culprit = usedByLocalStore ? "useLocalStore" : "useAsObservableSource";
    if ("production" !== process.env.NODE_ENV && usedByLocalStore) {
        var _a = __read(react_1.default.useState(current), 1), initialSource = _a[0];
        if ((initialSource !== undefined && current === undefined) ||
            (initialSource === undefined && current !== undefined)) {
            throw new Error("make sure you never pass `undefined` to " + culprit);
        }
    }
    if (usedByLocalStore && current === undefined) {
        return undefined;
    }
    if ("production" !== process.env.NODE_ENV && !utils_1.isPlainObject(current)) {
        throw new Error(culprit + " expects a plain object as " + (usedByLocalStore ? "second" : "first") + " argument");
    }
    var _b = __read(react_1.default.useState(function () { return mobx_1.observable(current, {}, { deep: false }); }), 1), res = _b[0];
    if ("production" !== process.env.NODE_ENV &&
        Object.keys(res).length !== Object.keys(current).length) {
        throw new Error("the shape of objects passed to " + culprit + " should be stable");
    }
    mobx_1.runInAction(function () {
        Object.assign(res, current);
    });
    return res;
}
exports.useAsObservableSourceInternal = useAsObservableSourceInternal;
function useAsObservableSource(current) {
    return useAsObservableSourceInternal(current, false);
}
exports.useAsObservableSource = useAsObservableSource;
