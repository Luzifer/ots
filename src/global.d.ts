export { }

declare global {
  interface Window {
    getTheme: () => string;
    getThemeFromStorage: () => string;
    maxSecretExpire: number;
    OTSCustomize: any;
    refreshTheme: () => void;
    setTheme: (theme: string) => void;
    useFormalLanguage: boolean;
    version: string;
  }
}
