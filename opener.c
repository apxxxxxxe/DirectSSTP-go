#include <stdio.h>
#include <string.h>
#include <windows.h>
#include "opener.h"

char* getFMO(void) {
  char* str;

  HANDLE h = OpenFileMapping(FILE_MAP_READ, FALSE, "Sakura");

  if (h == 0) {
    // printf("handle not found\n");
    CloseHandle(h);
    return(NULL);
  }

  int sharedMemSize = 1024 * 64;
  char* mapPtr = (char*)MapViewOfFile(h, FILE_MAP_READ, 0, 0, sharedMemSize);
  mapPtr += 4; // 最初の4バイトは必要ない情報

  if ( NULL != mapPtr ) {

    int i=0;
    while (mapPtr[i] != '\0') {
      i++;
    }

    str = malloc(i);
    strncpy(str, mapPtr, i);
    str[i] = '\0';

    // printf("%s", str);
    UnmapViewOfFile(mapPtr);
    CloseHandle(h);
    return(str);
  }
  else{
    // printf("failed\n");
    UnmapViewOfFile(mapPtr);
    CloseHandle(h);
    return(NULL);
  }
}
