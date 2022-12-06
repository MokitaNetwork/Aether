package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/mokitanetwork/aether/x/aethdist/client/cli"
	"github.com/mokitanetwork/aether/x/aethdist/client/rest"
)

// community-pool multi-spend proposal handler
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
