package irremote

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/xkamail/huberlink-platform/pkg/pgctx"
	"github.com/xkamail/huberlink-platform/pkg/snowid"
	"github.com/xkamail/huberlink-platform/pkg/uierr"
)

type Virtualer interface {
	// Properties get default properties
	Properties() Properties
	// PropertiesKeys validate properties
	PropertiesKeys() []string
}

// VirtualCategory is an enum that represents a virtual key category
type VirtualCategory uint

// Properties represent state of a virtual key (device)
type Properties map[string]string

const (
	VirtualCategoryOther VirtualCategory = iota
	VirtualCategoryTV
	VirtualCategoryAirConditioner
	VirtualCategoryLight
	VirtualCategoryFan
	VirtualCategorySpeaker
	VirtualCategoryProjector
	VirtualCategoryDVD
	VirtualCategoryWaterHeater
)

// IRRemote is a struct that represents a remote control
// that linked to a device.Device and home.Home
type IRRemote struct {
	ID        snowid.ID `json:"id"`
	DeviceID  snowid.ID `json:"deviceId"`
	HomeID    snowid.ID `json:"homeId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// VirtualKey is a struct that represents a remote control
// virtual remote control
type VirtualKey struct {
	ID         snowid.ID       `json:"id"`
	RemoteID   snowid.ID       `json:"remoteId"`
	Name       string          `json:"name"`
	Icon       string          `json:"icon"`
	Kind       VirtualCategory `json:"category"`
	Properties Properties      `json:"properties"`
	CreatedAt  time.Time       `json:"createdAt"`
	UpdatedAt  time.Time       `json:"updatedAt"`
}

func Find(ctx context.Context, deviceID snowid.ID) (*IRRemote, error) {
	rows, err := pgctx.Query(ctx, `
		select id, device_id, home_id, created_at, updated_at
		from device_ir_remotes
		where device_id = $1`,
		deviceID,
	)
	if err != nil {
		return nil, err
	}
	ir, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByPos[IRRemote])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrRemoteNotFound
	}
	if err != nil {
		return nil, err
	}
	return ir, nil
}

func ListVirtual(ctx context.Context, deviceID snowid.ID) ([]*VirtualKey, error) {
	return pgctx.Collect[VirtualKey](ctx, `
		select v.id, v.remote_id, v.name, v.icon, v.kind, v.properties, v.created_at, v.updated_at
		from device_ir_remote_virtual_keys v
		inner join device_ir_remotes d on d.id = v.remote_id
		where remote_id = $1`,
		deviceID,
	)
}

type CreateVirtualKeyParam struct {
	RemoteID snowid.ID       `json:"remoteId"`
	Name     string          `json:"name"`
	Category VirtualCategory `json:"kind"`
}

func (p *CreateVirtualKeyParam) Valid() error {
	p.Name = strings.TrimSpace(p.Name)
	if p.RemoteID == 0 {
		return uierr.Invalid("remote id", "remoteId is required")
	}
	if p.Name == "" {
		return uierr.Invalid("name", "name is required")
	}
	if len(p.Name) > 255 {
		return uierr.Invalid("name", "name is too long")
	}
	switch p.Category {
	case VirtualCategoryAirConditioner, VirtualCategoryTV, VirtualCategoryOther:
		return nil
	default:
		return uierr.Invalid("category", "category is invalid")
	}
}

func CreateVirtual(ctx context.Context, p *CreateVirtualKeyParam) (*VirtualKey, error) {
	if err := p.Valid(); err != nil {
		return nil, err
	}
	defaultProperties := make(Properties)
	switch p.Category {
	case VirtualCategoryTV:
		defaultProperties = TV{}.Properties()
	case VirtualCategoryAirConditioner:
		defaultProperties = Air{}.Properties()
	case VirtualCategoryOther:
		defaultProperties = Other{}.Properties()
	}
	now := time.Now()
	var virtualKeyID snowid.ID
	err := pgctx.QueryRow(ctx, `
		insert into device_ir_remote_virtual_keys 
			(id, remote_id, name, kind, icon, properties, created_at, updated_at) 
		values ($1,$2,$3,$4,$5,$6,$7,$8) returning id`,
		snowid.Gen(),
		p.RemoteID,
		p.Name,
		p.Category,
		"",
		defaultProperties,
		now,
		now,
	).Scan(&virtualKeyID)
	if err != nil {
		return nil, err
	}

	return FindVirtual(ctx, virtualKeyID)
}

func FindVirtual(ctx context.Context, virtualKeyID snowid.ID) (*VirtualKey, error) {
	rows, err := pgctx.Query(ctx, `
		select id, remote_id, name, kind, icon, properties, created_at, updated_at
		from device_ir_remote_virtual_keys 
		where id = $1`,
		virtualKeyID,
	)
	if err != nil {
		return nil, err
	}
	v, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByPos[VirtualKey])
	if err != nil {
		return nil, err
	}
	return v, nil
}
