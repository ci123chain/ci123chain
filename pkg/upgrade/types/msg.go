package types

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"strings"
)

// Plan specifies information about a planned upgrade and when it should occur
type Plan struct {
	// Sets the name for the upgrade. This name will be used by the upgraded version of the software to apply any
	// special "on-upgrade" commands during the first BeginBlock method after the upgrade is applied. It is also used
	// to detect whether a software version can handle a given upgrade. If no upgrade handler with this name has been
	// set in the software, it will be assumed that the software is out-of-date when the upgrade Time or Height
	// is reached and the software will exit.
	Name string `json:"name,omitempty"`

	// The time after which the upgrade must be performed.
	// Leave set to its zero value to use a pre-defined Height instead.
	//Time time.Time `json:"time,omitempty"`

	// The height at which the upgrade must be performed.
	// Only used if Time is not set.
	Height int64 `json:"height,omitempty"`

	// Any application specific upgrade info to be included on-chain
	// such as a git commit that validators could automatically upgrade to
	Info string `json:"info,omitempty"`
}

func (p Plan) String() string {
	due := p.DueAt()
	dueUp := strings.ToUpper(due[0:1]) + due[1:]
	return fmt.Sprintf(`Upgrade Plan
  Name: %s
  %s
  Info: %s`, p.Name, dueUp, p.Info)
}

// ValidateBasic does basic validation of a Plan
func (p Plan) ValidateBasic() error {
	if len(p.Name) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "name cannot be empty")
	}
	if p.Height < 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "height cannot be negative")
	}
	if p.Height == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "must set height")
	}

	return nil
}

// ShouldExecute returns true if the Plan is ready to execute given the current context
func (p Plan) ShouldExecute(ctx sdk.Context) bool {
	if p.Height > 0 {
		return p.Height <= ctx.BlockHeight()
	}
	return false
}

// DueAt is a string representation of when this plan is due to be executed
func (p Plan) DueAt() string {
	return fmt.Sprintf("height: %d", p.Height)
}

var _ sdk.Msg = &SoftwareUpgradeProposal{}

// Software Upgrade Proposals
type SoftwareUpgradeProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Plan        Plan   `json:"plan" yaml:"plan"`
	Proposer 	sdk.AccAddress `json:"proposer"`
}


func (msg *SoftwareUpgradeProposal) ValidateBasic() error {
	if msg.Proposer.Empty() {
		return sdkerrors.ErrInvalidParam
	}

	return msg.Plan.ValidateBasic()
}

func (msg *SoftwareUpgradeProposal) Bytes() []byte {
	bytes, err := UpgradeCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *SoftwareUpgradeProposal) Route() string {return RouterKey}
func (msg *SoftwareUpgradeProposal) MsgType() string {return "upgrade_proposal"}
func (msg *SoftwareUpgradeProposal) GetFromAddress() sdk.AccAddress { return msg.Proposer}



func NewSoftwareUpgradeProposal(proposer sdk.AccAddress, title, description string, plan Plan) *SoftwareUpgradeProposal {
	return &SoftwareUpgradeProposal{title, description, plan, proposer}
}


// Cancel Software Upgrade Proposals
type CancelSoftwareUpgradeProposal struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Proposer 	sdk.AccAddress `json:"proposer"`
}


func (msg *CancelSoftwareUpgradeProposal) ValidateBasic() error {
	if msg.Proposer.Empty() {
		return sdkerrors.ErrInvalidParam
	}

	return nil
}

func (msg *CancelSoftwareUpgradeProposal) Bytes() []byte {
	bytes, err := UpgradeCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *CancelSoftwareUpgradeProposal) Route() string {return RouterKey}
func (msg *CancelSoftwareUpgradeProposal) MsgType() string {return "upgrade_proposal"}
func (msg *CancelSoftwareUpgradeProposal) GetFromAddress() sdk.AccAddress { return msg.Proposer}


func NewCancelSoftwareUpgradeProposal(proposer sdk.AccAddress) *CancelSoftwareUpgradeProposal {
	return &CancelSoftwareUpgradeProposal{Proposer: proposer}
}

