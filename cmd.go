package hermeti

import "context"

type Command func(context.Context, Env, []string) []string
