package cmd

import (
	"context"
	"retoBack/internal/cmd/routes" // 👈 importa tu paquete de rutas

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()

			// Registrar todas las rutas aquí
			routes.Register(s)

			// Levantar el servidor
			s.Run()
			return nil
		},
	}
)
