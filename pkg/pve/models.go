package pve

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Describes proxmox API JSON objects for external use

/* Descirbed models paths


 */

type VMID uint64

func (id *VMID) UnmarshalJSON(data []byte) error {

	var num uint64
	if err := json.Unmarshal(data, &num); err == nil {
		*id = VMID(num)
		return nil
	}

	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("VMID: cannot parse as string or number: %w", err)
	}

	parsed, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return fmt.Errorf("VMID: cannot parse string to uint64: %w", err)
	}

	*id = VMID(parsed)
	return nil
}

func (id VMID) String() string {
	return strconv.FormatUint(uint64(id), 10)
}

// ResponseWrapper - proxmox JSON object response base wrap from hots to client in successfull request
type ResponseWrapper[T any] struct {
	Data    T      `json:"data"`
	Message string `json:"message,omitempty"`
}

type NodeBaseInfo struct {
	Maxdisk        uint64  `json:"maxdisk"`
	CPU            float64 `json:"cpu"`
	Disk           uint64  `json:"disk"`
	ID             string  `json:"id"`
	Mem            uint64  `json:"mem"`
	Level          string  `json:"level"`
	Node           string  `json:"node"`
	Type           string  `json:"type"`
	Status         string  `json:"status"`
	SslFingerprint string  `json:"ssl_fingerprint"`
	Maxcpu         uint64  `json:"maxcpu"`
	Uptime         int64   `json:"uptime"`
	Maxmem         uint64  `json:"maxmem"`
}

func (i NodeBaseInfo) NodeName() string {
	return i.Node
}

type Nodes ResponseWrapper[[]NodeBaseInfo]
type Node ResponseWrapper[NodeBaseInfo]

type (
	CpuInfo struct {
		Cpus    uint32 `json:"cpus"`
		Flags   string `json:"flags"`
		Mhz     string `json:"mhz"`
		Cores   uint32 `json:"cores"`
		Model   string `json:"model"`
		UserHz  uint32 `json:"user_hz"`
		Hvm     string `json:"hvm"`
		Sockets uint32 `json:"sockets"`
	}

	BootInfo struct {
		SecureBoot uint64 `json:"secureboot"`
		Mode       string `json:"mode"`
	}

	Memory struct {
		Used  uint64 `json:"used"`
		Free  uint64 `json:"free"`
		Total uint64 `json:"total"`
	}

	CurrentKernel struct {
		Machine string `json:"machine"`
		Sysname string `json:"sysname"`
		Release string `json:"release"`
		Version string `json:"version"`
	}

	Swap struct {
		Used  uint64 `json:"used"`
		Free  uint64 `json:"free"`
		Total uint64 `json:"total"`
	}

	RootFS struct {
		Used  uint64 `json:"used"`
		Free  uint64 `json:"free"`
		Total uint64 `json:"total"`
		Avail uint64 `json:"avail"`
	}

	LoadAvg [3]string

	NodeStat struct {
		Cpu           float64 `json:"cpu"`
		CpuInfo       `json:"cpuinfo"`
		Memory        `json:"memory"`
		Swap          `json:"swap"`
		RootFS        `json:"rootfs"`
		LoadAvg       `json:"loadavg"`
		Uptime        uint64 `json:"uptime"`
		BootInfo      `json:"boot-info"`
		CurrentKernel `json:"current-kernel"`
	}
)

func (l LoadAvg) Load5() float64 {
	value, err := strconv.ParseFloat(l[0], 64)
	if err != nil {
		return 0.0
	}
	return value
}

func (l LoadAvg) Load10() float64 {
	value, err := strconv.ParseFloat(l[1], 64)
	if err != nil {
		return 0.0
	}
	return value
}

func (l LoadAvg) Load15() float64 {
	value, err := strconv.ParseFloat(l[2], 64)
	if err != nil {
		return 0.0
	}
	return value
}

type NodeStatus ResponseWrapper[NodeStat]

type InfoLXC struct {
	Vmid      VMID    `json:"vmid"`
	Pid       uint64  `json:"pid"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	CPU       float64 `json:"cpu"`
	Cpus      uint64  `json:"cpus"`
	Uptime    uint64  `json:"uptime"`
	Status    string  `json:"status"`
	Tags      string  `json:"tags"`
	Mem       uint64  `json:"mem"`
	Maxmem    uint64  `json:"maxmem"`
	Swap      uint64  `json:"swap"`
	Maxswap   uint64  `json:"maxswap"`
	Disk      uint64  `json:"disk"`
	Maxdisk   uint64  `json:"maxdisk"`
	Diskwrite uint64  `json:"diskwrite"`
	Diskread  uint64  `json:"diskread"`
	Netin     uint64  `json:"netin"`
	Netout    uint64  `json:"netout"`
}

type InfoQemu struct {
	Vmid      VMID    `json:"vmid"`
	Pid       uint64  `json:"pid"`
	Name      string  `json:"name"`
	CPU       float64 `json:"cpu"`
	Cpus      uint64  `json:"cpus"`
	Uptime    uint64  `json:"uptime"`
	Status    string  `json:"status"`
	Mem       uint64  `json:"mem"`
	Maxmem    uint64  `json:"maxmem"`
	Disk      uint64  `json:"disk"`
	Maxdisk   int64   `json:"maxdisk"`
	Diskwrite uint64  `json:"diskwrite"`
	Diskread  uint64  `json:"diskread"`
	Netin     uint64  `json:"netin"`
	Netout    uint64  `json:"netout"`
}

type (
	LxcData  ResponseWrapper[[]InfoLXC]
	QemuData ResponseWrapper[[]InfoQemu]
)

type (
	NodeStorage struct {
		Used         uint64  `json:"used"`
		Type         string  `json:"type"`
		UsedFraction float64 `json:"used_fraction"`
		Content      string  `json:"content"`
		Storage      string  `json:"storage"`
		Avail        uint64  `json:"avail"`
		Total        uint64  `json:"total"`
		Enabled      uint64  `json:"enabled"`
		Shared       uint64  `json:"shared"`
		Active       uint64  `json:"active"`
	}

	NodeStorageList ResponseWrapper[[]NodeStorage]
)

type (
	IfaceNetstat struct {
		Vmid VMID   `json:"vmid"`
		Dev  string `json:"dev"`
		In   string `json:"in"`
		Out  string `json:"out"`
	}

	NodeIfaceNetstatList ResponseWrapper[[]IfaceNetstat]
)

func (s IfaceNetstat) TX() float64 {
	value, err := strconv.ParseFloat(s.Out, 64)
	if err != nil {
		return 0.0
	}
	return value
}

func (s IfaceNetstat) RX() float64 {
	value, err := strconv.ParseFloat(s.In, 64)
	if err != nil {
		return 0.0
	}
	return value
}
