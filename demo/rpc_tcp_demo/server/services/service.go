package services

type Args struct {
	X int
	Y int
}

type ServiceA struct {
}

func (c *ServiceA) Add(args *Args, reply *int) error {
	*reply = args.X + args.Y
	return nil
}
