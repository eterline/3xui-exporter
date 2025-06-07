package pve

import (
	"context"
	"fmt"
)

type ProxmoxClient struct {
	Req *proxmoxRequests
}

func NewProxmoxClient(api string, tokenID, token, caFile string) (*ProxmoxClient, error) {

	requests, err := newProxmoxRequests(api, tokenID, token, caFile)
	if err != nil {
		return nil, fmt.Errorf("failed to init client: %w", err)
	}

	client := &ProxmoxClient{
		Req: requests,
	}

	return client, nil
}

/*
Nodes - request for node list in proxmox

	detail: `https://pve.proxmox.com/pve-docs/api-viewer/#/nodes`
*/
func (client *ProxmoxClient) Nodes(ctx context.Context) (Nodes, error) {
	req := client.Req.request(ctx, "nodes")
	data := Nodes{}

	code, err := req.get()
	if err != nil {
		return data, err
	}

	if code > 299 || code < 199 {
		return data, fmt.Errorf("bad status code: %d", code)
	}

	if err := req.resolve(&data); err != nil {
		return data, err
	}

	return data, nil
}

type NodeNamer interface {
	NodeName() string
}

func (client *ProxmoxClient) NodeStats(ctx context.Context, n NodeNamer) (NodeStatus, error) {
	req := client.Req.request(ctx, "nodes", n.NodeName(), "status")
	data := NodeStatus{}

	code, err := req.get()
	if err != nil {
		return data, err
	}

	if code > 299 || code < 199 {
		return data, fmt.Errorf("bad status code: %d", code)
	}

	if err := req.resolve(&data); err != nil {
		return data, err
	}

	return data, nil
}

func (client *ProxmoxClient) LxcStats(ctx context.Context, n NodeNamer) (LxcData, error) {
	req := client.Req.request(ctx, "nodes", n.NodeName(), "lxc")
	data := LxcData{}

	code, err := req.get()
	if err != nil {
		return data, err
	}

	if code > 299 || code < 199 {
		return data, fmt.Errorf("bad status code: %d", code)
	}

	if err := req.resolve(&data); err != nil {
		return data, err
	}

	return data, nil
}

func (client *ProxmoxClient) QemuStats(ctx context.Context, n NodeNamer) (QemuData, error) {
	req := client.Req.request(ctx, "nodes", n.NodeName(), "qemu")
	data := QemuData{}

	code, err := req.get()
	if err != nil {
		return data, err
	}

	if code > 299 || code < 199 {
		return data, fmt.Errorf("bad status code: %d", code)
	}

	if err := req.resolve(&data); err != nil {
		return data, err
	}

	return data, nil
}

func (client *ProxmoxClient) NodeStorages(ctx context.Context, n NodeNamer) (NodeStorageList, error) {
	req := client.Req.request(ctx, "nodes", n.NodeName(), "storage")
	data := NodeStorageList{}

	code, err := req.get()
	if err != nil {
		return data, err
	}

	if code > 299 || code < 199 {
		return data, fmt.Errorf("bad status code: %d", code)
	}

	if err := req.resolve(&data); err != nil {
		return data, err
	}

	return data, nil
}

func (client *ProxmoxClient) NodeNetstat(ctx context.Context, n NodeNamer) (NodeIfaceNetstatList, error) {
	req := client.Req.request(ctx, "nodes", n.NodeName(), "netstat")
	data := NodeIfaceNetstatList{}

	code, err := req.get()
	if err != nil {
		return data, err
	}

	if code > 299 || code < 199 {
		return data, fmt.Errorf("bad status code: %d", code)
	}

	if err := req.resolve(&data); err != nil {
		return data, err
	}

	return data, nil
}
