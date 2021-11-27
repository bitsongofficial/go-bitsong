import { configure } from "mobx";
import { getGlobal, getSymbol } from "./utils";
var observerBatchingConfiguredSymbol = getSymbol("observerBatching");
export function defaultNoopBatch(callback) {
    callback();
}
export function observerBatching(reactionScheduler) {
    if (!reactionScheduler) {
        reactionScheduler = defaultNoopBatch;
        if ("production" !== process.env.NODE_ENV) {
            console.warn("[MobX] Failed to get unstable_batched updates from react-dom / react-native");
        }
    }
    configure({ reactionScheduler: reactionScheduler });
    getGlobal()[observerBatchingConfiguredSymbol] = true;
}
export var isObserverBatched = function () { return !!getGlobal()[observerBatchingConfiguredSymbol]; };
