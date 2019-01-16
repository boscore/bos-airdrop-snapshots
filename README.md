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

### Total

- Create Accounts Usage: `315475.0000 BOS`
- Airdrop Usage: `50088508.9033 BOS`
- Total Usage: `50403983.9033 BOS`
- Total Accounts: `630950`
