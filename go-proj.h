#ifndef GO_PROJ_H
#define GO_PROJ_H

#include <proj.h>

#if PROJ_VERSION_MAJOR < 8
const char *proj_context_errno_string(PJ_CONTEXT *ctx, int err);
#endif

#endif