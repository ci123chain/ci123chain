package collactor

import "time"

func Link(src, dst *Chain, retries uint64, to time.Duration) error {
	// create clients if they aren't already created
	modified, err := src.CreateClients(dst)
	if modified {
		//if err := overWriteConfig(config); err != nil {
		//	return err
		//}
	}
	if err != nil {
		return err
	}


	// create connection if it isn't already created
	modified, err = src.CreateOpenConnections(dst, retries, to)
	if modified {
		//if err := overWriteConfig(config); err != nil {
		//	return err
		//}
	}
	if err != nil {
		return err
	}

	// create channel if it isn't already created
	modified, err = src.CreateOpenChannels(dst, retries, to)
	if modified {
		//if err := overWriteConfig(config); err != nil {
		//	return err
		//}
	}
	return err
}
