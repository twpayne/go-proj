#include "go-proj.h"

#if PROJ_VERSION_MAJOR < 8
const char *proj_context_errno_string(PJ_CONTEXT *ctx, int err) {
    return proj_errno_string(err);
}
#endif