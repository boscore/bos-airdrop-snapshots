# README

Snapshots for the BOS Mainnet airdrop for EOS Mainnet accounts.

### Snapshot Info

- Snapshot we took on EOS Mainnet at `2019-Jan-01, 09:50 UTC+0`, about block height `#35072000`.
- `accounts_info_bos_snapshot.airdrop.normal.csv` Airdrop file for the BOS Mainnet taken form the EOS Mainnet snapshot for non-msig accounts.
- `accounts_info_bos_snapshot.airdrop.msig.json` Airdrop file for the BOS Mainnet taken form the EOS Mainnet snapshot for msig accounts.

### Strategy

Our strategy for generating the BOS Mainnet airdrop snapshot is as follows:

1. For non-msig accounts, calculate the balance of an EOS Mainnet account on BOS Mainnet, `20.0000 EOS ` on EOS Mainnet can get `1.0000 BOS` on BOS Mainnet and the accounts whose balance is smaller than `0.5000 BOS` will be given at least `0.5000 BOS`;
2. For msig accounts, the balances are calculated as the same as non-msig accounts, and we do an auth mapping of it's permissions.
3. Generate a random name for EOS Mainnet accounts;
4. Accounts are created with: `0.3000 BOS` ram, `0.1000 BOS` NET  and `0.1000 BOS` CPU on BOS Mainnet;
5. Generate the snapshot for airdrop.
6. Every account will keep 10 BOS at most, and the left will be delegated to CPU and NET 

### Total

- Create Accounts Usage: `315475 BOS`
- Airdrop Usage: `50088509 BOS`
- Total Usage: `50403984 BOS`
- Total Accounts: `630949`


### Inactive Airdrop BOS Burning Proposal

BOSCore's main network started with BOS token airdrops on more than 630,000 accounts, of which there are still a lot of accounts that have never been activated before (in this proposal, "inactive accounts" means `auth_sequence=0||auth_sequence=2`). In order to better promote the follow-up development of the BOSCore mainnet and community, this proposal proposes to burn the BOS tokens of the inactive accounts and will be implemented within 2 weeks after the proposal is approved.

The snapshot for burning time is `2019-11-27 03:26:27 UTC-0`，block height is `54,171,828`. The inactive accounts file is [unactive_airdrop_accounts.csv](./unactive_airdrop_accounts.csv).


