"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var mobx_1 = require("mobx");
var utils_1 = require("./utils");
var observerBatchingConfiguredSymbol = utils_1.getSymbol("observerBatching");
function defaultNoopBatch(callback) {
    callback();
}
exports.defaultNoopBatch = defaultNoopBatch;
function observerBatching(reactionScheduler) {
    if (!reactionScheduler) {
        reactionScheduler = defaultNoopBatch;
        if ("production" !== process.env.NODE_ENV) {
            console.warn("[MobX] Failed to get unstable_batched updates from react-dom / react-native");
        }
    }
    mobx_1.configure({ reactionScheduler: reactionScheduler });
    utils_1.getGlobal()[observerBatchingConfiguredSymbol] = true;
}
exports.observerBatching = observerBatching;
exports.isObserverBatched = function () { return !!utils_1.getGlobal()[observerBatchingConfiguredSymbol]; };
