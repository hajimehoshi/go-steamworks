package steamworks

// extern void warningMessageGoHook(int severity, void * v);
//
// void __cdecl SteamAPIDebugTextGlobalHook(int nSeverity, const char *pchDebugText)
// {
//   warningMessageGoHook(nSeverity, (void *)pchDebugText);
// }
import "C"
