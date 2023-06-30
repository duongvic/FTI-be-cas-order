package migrations

import (
	"casorder/db/models"
	"casorder/utils/types"
)

var UnitData = []models.Unit {
	{
		Name: "Kilobyte",
		Code: "KB",
		Description: "Kilobyte",
	}, 
	{
		Name: "Megabyte",
		Code: "MB",
		Description: "Megabyte",
	}, 
	{
		Name: "Gigabyte",
		Code: "GB",
		Description: "Gigabyte",
	}, 
	{
		Name: "Terabyte",
		Code: "TB",
		Description: "Terabyte",
	}, 
	{
		Name: "Petabyte",
		Code: "PB",
		Description: "Petabyte",
	}, 
	{
		Name: "vCPU",
		Code: "vCPU",
		Description: "Virtual CPU",
	}, 
	{
		Name: "License",
		Code: "License",
		Description: "License. Take Window License for an example",
	}, 
	{
		Name: "IP",
		Code: "IP",
		Description: "",
	}, 
	{
		Name: "Mbps",
		Code: "Mbps",
		Description: "",
	}, 
	{
		Name: "Day",
		Code: "Day",
		Description: "",
	}, 
	{
		Name: "Month",
		Code: "Month",
		Description: "",
	}, 
	{
		Name: "Year",
		Code: "Year",
		Description: "",
	},
	{
		Name: "Count",
		Code: "Count",
		Description: "",
	},
}

var ProductData = []models.Product{
	{
		Name:           "CPU",
		CN:             "vcpu",
		Type:           "PROCESSOR",
		Description:    "vCPU (Intel Xeon)",
		InitFee:        0,
		MaintenanceFee: 100,
		IsBase:         true,
		UnitID:         6,
	},
	{
		Name:           "MEMORY",
		CN:             "ram",
		Type:           "MEMORY",
		Description:    "GB",
		InitFee:        0,
		MaintenanceFee: 100,
		IsBase:         true,
		UnitID:         3,
	},
	{
		Name:           "DISK",
		CN:             "disk",
		Type:           "STORAGE",
		Description:    "GB",
		InitFee:        0,
		MaintenanceFee: 100,
		IsBase:         true,
		UnitID:         3,
	},
	{
		Name:           "IP",
		CN:             "ip",
		Type:           "NET",
		Description:    "",
		InitFee:        0,
		MaintenanceFee: 100,
		IsBase:         true,
		UnitID:         8,
	},
	{
		Name:           "NET",
		CN:             "net",
		Type:           "NET",
		Description:    "",
		InitFee:        0,
		MaintenanceFee: 100,
		IsBase:         true,
		UnitID:         9,
	},
	{
		Name:           "Ubuntu 12.04",
		CN:             "image",
		Type:           "OS",
		Description:    "",
		InitFee:        0,
		MaintenanceFee: 100,
		IsBase:         true,
		UnitID:         7,
		Data: types.MapToEncodedJSON(map[string]interface{}{
			"arch":         "x64",
			"type":         "ubuntu",
			"version":      "12.04 LTS",
			"platform":     "server",
			"backend_name": "ubuntu-12.04LTS-x64-server",
		}),
	},
	{
		Name:           "Ubuntu 16.04",
		CN:             "image",
		Type:           "OS",
		Description:    "",
		InitFee:        0,
		MaintenanceFee: 100,
		IsBase:         true,
		UnitID:         7,
		Data: types.MapToEncodedJSON(map[string]interface{}{
			"arch":         "x64",
			"type":         "ubuntu",
			"version":      "16.04 LTS",
			"platform":     "server",
			"backend_name": "ubuntu-16.04LTS-x64-server",
		}),
	},
	{
		Name:           "Windows 10",
		CN:             "image",
		Type:           "OS",
		Description:    "",
		InitFee:        0,
		MaintenanceFee: 100,
		IsBase:         true,
		UnitID:         7,
		Data: types.MapToEncodedJSON(map[string]interface{}{
			"arch":         "x64",
			"type":         "window",
			"version":      "10",
			"platform":     "server",
			"backend_name": "window-10Pro-x64-server", // {type}-{version}-{arch}-{platform}
		}),
	},
	{
		Name:           "Snapshot",
		CN:             "snapshot",
		Type:           "STORAGE",
		Description:    "",
		InitFee:        0,
		MaintenanceFee: 100,
		IsBase:         false,
		UnitID:         3,
	},
	{
		Name:           "Backup",
		CN:             "backup",
		Type:           "STORAGE",
		Description:    "",
		InitFee:        0,
		MaintenanceFee: 100,
		IsBase:         false,
		UnitID:         3,
	},
	{
		Name:           "ROOT_DISK",
		CN:             "root_disk",
		Type:           "STORAGE",
		Description:    "GB",
		InitFee:        0,
		MaintenanceFee: 100,
		IsBase:         true,
		UnitID:         3,
	},
	{
		Name:           "DATA_DISK",
		CN:             "data_disk",
		Type:           "STORAGE",
		Description:    "GB",
		InitFee:        0,
		MaintenanceFee: 100,
		IsBase:         true,
		UnitID:         3,
	},
	{
		Name:           "VPN",
		CN:             "vpn",
		Type:           "NET",
		Description:    "Count",
		InitFee:        0,
		MaintenanceFee: 100,
		IsBase:         false,
		UnitID:         13,
	},
	{
		Name:           "LOAD_BALANCER",
		CN:             "lb",
		Type:           "NET",
		Description:    "Count",
		InitFee:        0,
		MaintenanceFee: 100,
		IsBase:         false,
		UnitID:         13,
	},
	{
		Name:           "GPU",
		CN:             "gpu",
		Type:           "PROCESSOR",
		Description:    "Count",
		InitFee:        0,
		MaintenanceFee: 100,
		IsBase:         false,
		UnitID:         13,
	},
}

var PackageData = []models.Package{
	{
		Name:           "Custom",
		Type:           "STANDARD",
		UnitID:         11,
		TrialTime:      0,
		AllowCustom:    true,
		InitFee:        0,
		MaintenanceFee: 100,
	},
	{
		Name:           "Small",
		Type:           "STANDARD",
		UnitID:         11,
		TrialTime:      14,
		InitFee:        0,
		MaintenanceFee: 100,
	},
	{
		Name:           "Huge",
		Type:           "STANDARD",
		UnitID:         11,
		TrialTime:      14,
		InitFee:        0,
		MaintenanceFee: 100,
	},
}

var RegionData = []models.Region {
	{
		Description: "HN",
		Name: "HN",
		Status: true,
	},
	{
		Description: "HCM",
		Name: "HCM",
		Status: true,
	},
	{
		Description: "DN",
		Name: "DN",
		Status: true,
	},
}


var PackageProductData = []models.PackageProduct {
	{
		PackageID: 1,
		ProductID: 1,
		UnitID: 6,
		Quantity: 0,
	},
	{
		PackageID: 1,
		ProductID: 2,
		UnitID: 3,
		Quantity: 0,
	},
	{
		PackageID: 1,
		ProductID: 3,
		UnitID: 3,
		Quantity: 0,
	},
	{
		PackageID: 1,
		ProductID: 4,
		UnitID: 6,
		Quantity: 0,
	},
	{
		PackageID: 1,
		ProductID: 5,
		UnitID: 3,
		Quantity: 0,
	},
	{
		PackageID: 1,
		ProductID: 6,
		UnitID: 10,
		Quantity: 0,
	},
	{
		PackageID: 1,
		ProductID: 7,
		UnitID: 11,
		Quantity: 0,
	},
	{
		PackageID: 2,
		ProductID: 1,
		UnitID: 6,
		Quantity: 100,
	},
	{
		PackageID: 2,
		ProductID: 2,
		UnitID: 3,
		Quantity: 100,
	},
	{
		PackageID: 2,
		ProductID: 3,
		UnitID: 3,
		Quantity: 100,
	},
}

