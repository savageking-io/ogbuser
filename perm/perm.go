package perm

import (
	"fmt"
	"github.com/savageking-io/ogbuser/schema"
	log "github.com/sirupsen/logrus"
)

const (
	DomainOwn    string = "own"
	DomainParty  string = "party"
	DomainGuild  string = "guild"
	DomainGlobal string = "global"
)

type Permission struct {
	Name   string
	Read   int32
	Write  int32
	Delete int32
	Domain string
	raw    schema.GroupPermissionSchema
}

type Perm struct {
	own         map[string]Permission
	party       map[string]Permission
	guild       map[string]Permission
	global      map[string]Permission
	ownArray    []*Permission
	partyArray  []*Permission
	guildArray  []*Permission
	globalArray []*Permission
}

func NewPerm() *Perm {
	return &Perm{
		own:    make(map[string]Permission),
		party:  make(map[string]Permission),
		guild:  make(map[string]Permission),
		global: make(map[string]Permission),
	}
}

func (p *Perm) Add(domain string, perm Permission) {
	if domain == DomainOwn {
		p.AddOwn(perm)
		return
	}
	if domain == DomainParty {
		p.AddParty(perm)
		return
	}
	if domain == DomainGuild {
		p.AddGuild(perm)
		return
	}
	if domain == DomainGlobal {
		p.AddGlobal(perm)
		return
	}
	log.Errorf("Perm::Add() Unknown domain: %s", domain)
}

func (p *Perm) Populate(rawPerm *schema.GroupPermissionSchema) error {
	if rawPerm == nil {
		return fmt.Errorf("rawPerm is nil")
	}

	if rawPerm.Permission == "" {
		return fmt.Errorf("permission name is empty")
	}

	if rawPerm.Domain == "" {
		return fmt.Errorf("permission domain is empty")
	}

	perm := Permission{}
	perm.raw = *rawPerm
	perm.Name = perm.raw.Permission
	perm.Domain = perm.raw.Domain
	perm.Read = schema.BoolToInt32(perm.raw.Read)
	perm.Write = schema.BoolToInt32(perm.raw.Write)
	perm.Delete = schema.BoolToInt32(perm.raw.Delete)
	p.Add(perm.raw.Domain, perm)
	return nil
}

func (p *Perm) AddOwn(perm Permission) {
	p.own[perm.Name] = perm
	p.ownArray = append(p.ownArray, &perm)
}

func (p *Perm) AddParty(perm Permission) {
	p.party[perm.Name] = perm
	p.partyArray = append(p.partyArray, &perm)
}

func (p *Perm) AddGuild(perm Permission) {
	p.guild[perm.Name] = perm
	p.guildArray = append(p.guildArray, &perm)
}

func (p *Perm) AddGlobal(perm Permission) {
	p.global[perm.Name] = perm
	p.globalArray = append(p.globalArray, &perm)
}

func (p *Perm) Get(domain string) []*Permission {
	if domain == DomainOwn {
		return p.GetOwn()
	}
	if domain == DomainParty {
		return p.GetParty()
	}
	if domain == DomainGuild {
		return p.GetGuild()
	}
	if domain == DomainGlobal {
		return p.GetGlobal()
	}
	log.Errorf("Perm::Get() Unknown domain: %s", domain)
	return nil
}

func (p *Perm) GetOwn() []*Permission {
	return p.ownArray
}

func (p *Perm) GetParty() []*Permission {
	return p.partyArray
}

func (p *Perm) GetGuild() []*Permission {
	return p.guildArray
}

func (p *Perm) GetGlobal() []*Permission {
	return p.globalArray
}

func (p *Perm) GetPermission(domain string, permission string) *Permission {
	if domain == DomainOwn {
		return p.GetPermOwn(permission)
	}
	if domain == DomainParty {
		return p.GetPermParty(permission)
	}
	if domain == DomainGuild {
		return p.GetPermGuild(permission)
	}
	if domain == DomainGlobal {
		return p.GetPermGlobal(permission)
	}
	log.Errorf("Perm::Get() Unknown domain: %s", domain)
	return nil
}

func (p *Perm) GetPermOwn(permission string) *Permission {
	perm, ok := p.own[permission]
	if !ok {
		return &Permission{}
	}
	return &perm
}

func (p *Perm) GetPermParty(permission string) *Permission {
	perm, ok := p.party[permission]
	if !ok {
		return &Permission{}
	}
	return &perm
}

func (p *Perm) GetPermGuild(permission string) *Permission {
	perm, ok := p.guild[permission]
	if !ok {
		return &Permission{}
	}
	return &perm
}

func (p *Perm) GetPermGlobal(permission string) *Permission {
	perm, ok := p.global[permission]
	if !ok {
		return &Permission{}
	}
	return &perm
}

func (p *Perm) Count() int {
	return len(p.own) + len(p.party) + len(p.guild) + len(p.global)
}
