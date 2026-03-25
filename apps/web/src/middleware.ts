import createMiddleware from "next-intl/middleware";
import { routing } from "./i18n/routing";

export default createMiddleware(routing);

export const config = {
  // Exclude static files (manifest, icons, etc.) so next-intl does not prefix
  // them with a locale — otherwise /manifest.json becomes /en/manifest.json (404).
  matcher: ["/((?!api|_next|_vercel|brand|favicon\\.ico|.*\\..*).*)"],
};
