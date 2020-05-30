package v1

import (
	"context"
	"summa-auth-api/prisma-client"
)

var Client *prisma.Client
var Context = context.Background()

func LoadPrisma(){
	Client = prisma.New(nil)
}
