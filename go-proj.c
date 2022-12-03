#include "go-proj.h"

#if PROJ_VERSION_MAJOR < 8
const char *proj_context_errno_string(PJ_CONTEXT *ctx, int err) {
  return proj_errno_string(err);
}
#endif

#if PROJ_VERSION_MAJOR < 8 ||                                                  \
    (PROJ_VERSION_MAJOR == 8 && PROJ_VERSION_MINOR < 2)
int proj_trans_bounds(PJ_CONTEXT *context, PJ *P, PJ_DIRECTION direction,
                      double xmin, double ymin, double xmax, double ymax,
                      double *out_xmin, double *out_ymin, double *out_xmax,
                      double *out_ymax, int densify_pts) {
  return 1;
}
#endif

#if PROJ_VERSION_MAJOR < 9 ||                                                  \
    (PROJ_VERSION_MAJOR == 9 && PROJ_VERSION_MINOR < 1)
PJ *proj_trans_get_last_used_operation(PJ *P) { return NULL; }
#endif