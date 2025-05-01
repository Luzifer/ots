export { }

declare global {
  interface Window {
    getTheme: () => string;
    getThemeFromStorage: () => string;
    maxSecretExpire: number;
    OTSCustomize: any;
    refreshTheme: () => void;
    setTheme: (string) => void;
    useFormalLanguage: boolean;
    version: string;
  }
}
