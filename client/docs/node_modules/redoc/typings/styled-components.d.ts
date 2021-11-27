import * as styledComponents from 'styled-components';
import { ResolvedThemeInterface } from './theme';
export { ResolvedThemeInterface };
declare const styled: styledComponents.ThemedStyledInterface<ResolvedThemeInterface>, css: styledComponents.ThemedCssFunction<ResolvedThemeInterface>, createGlobalStyle: <P extends object = {}>(first: styledComponents.CSSObject | TemplateStringsArray | styledComponents.InterpolationFunction<styledComponents.ThemedStyledProps<P, ResolvedThemeInterface>>, ...interpolations: styledComponents.Interpolation<styledComponents.ThemedStyledProps<P, ResolvedThemeInterface>>[]) => styledComponents.GlobalStyleComponent<P, ResolvedThemeInterface>, keyframes: {
    (strings: TemplateStringsArray | styledComponents.CSSKeyframes, ...interpolations: styledComponents.SimpleInterpolation[]): styledComponents.Keyframes;
    (strings: string[] | TemplateStringsArray, ...interpolations: styledComponents.SimpleInterpolation[]): styledComponents.Keyframes;
}, ThemeProvider: styledComponents.BaseThemeProviderComponent<ResolvedThemeInterface, ResolvedThemeInterface>;
export declare const media: {
    lessThan(breakpoint: any, print?: boolean | undefined, extra?: string | undefined): (...args: any[]) => styledComponents.FlattenInterpolation<styledComponents.ThemedStyledProps<object, ResolvedThemeInterface>>;
    greaterThan(breakpoint: any): (...args: any[]) => styledComponents.FlattenInterpolation<styledComponents.ThemedStyledProps<object, ResolvedThemeInterface>>;
    between(firstBreakpoint: any, secondBreakpoint: any): (...args: any[]) => styledComponents.FlattenInterpolation<styledComponents.ThemedStyledProps<object, ResolvedThemeInterface>>;
};
export { css, createGlobalStyle, keyframes, ThemeProvider };
export default styled;
export declare function extensionsHook(styledName: string): (props: any) => any;
