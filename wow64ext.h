typedef unsigned long long DWORD64;
void WOW64Init();
DWORD64 GetProcessModuleHandle64(int hProcess, const wchar_t* lpModuleName, wchar_t *fullname); 
int ReadProcessMemory64(int hProcess, DWORD64 lpBaseAddress, void* lpBuffer, size_t nSize, size_t *lpNumberOfBytesRead);

