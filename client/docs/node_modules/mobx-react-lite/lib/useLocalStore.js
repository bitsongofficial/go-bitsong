"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
var mobx_1 = require("mobx");
var react_1 = __importDefault(require("react"));
var useAsObservableSource_1 = require("./useAsObservableSource");
var utils_1 = require("./utils");
function useLocalStore(initializer, current) {
    var source = useAsObservableSource_1.useAsObservableSourceInternal(current, true);
    return react_1.default.useState(function () {
        var local = mobx_1.observable(initializer(source));
        if (utils_1.isPlainObject(local)) {
            mobx_1.runInAction(function () {
                Object.keys(local).forEach(function (key) {
                    var value = local[key];
                    if (typeof value === "function") {
                        // @ts-ignore No idea why ts2536 is popping out here
                        local[key] = wrapInTransaction(value, local);
                    }
                });
            });
        }
        return local;
    })[0];
}
exports.useLocalStore = useLocalStore;
// tslint:disable-next-line: ban-types
function wrapInTransaction(fn, context) {
    return function () {
        var args = [];
        for (var _i = 0; _i < arguments.length; _i++) {
            args[_i] = arguments[_i];
        }
        return mobx_1.transaction(function () { return fn.apply(context, args); });
    };
}
