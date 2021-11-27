import { observable, runInAction, transaction } from "mobx";
import React from "react";
import { useAsObservableSourceInternal } from "./useAsObservableSource";
import { isPlainObject } from "./utils";
export function useLocalStore(initializer, current) {
    var source = useAsObservableSourceInternal(current, true);
    return React.useState(function () {
        var local = observable(initializer(source));
        if (isPlainObject(local)) {
            runInAction(function () {
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
// tslint:disable-next-line: ban-types
function wrapInTransaction(fn, context) {
    return function () {
        var args = [];
        for (var _i = 0; _i < arguments.length; _i++) {
            args[_i] = arguments[_i];
        }
        return transaction(function () { return fn.apply(context, args); });
    };
}
