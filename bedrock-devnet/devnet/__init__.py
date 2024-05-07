import argparse
import logging
import os
import jwt
import subprocess
import requests
import json
import socket
import calendar
import datetime
import time
import shutil
import http.client
from multiprocessing import Process, Queue
import concurrent.futures
from collections import namedtuple
# from hdwallet import BIP44HDWallet
# from hdwallet.cryptocurrencies import EthereumMainnet
# from hdwallet.derivations import BIP44Derivation
# from hdwallet.utils import generate_mnemonic
# from typing import Optional


import devnet.log_setup

pjoin = os.path.join

JWT_EXPIRATION_SECONDS = 3600  # 1 hour expiration time for the token

parser = argparse.ArgumentParser(description='Bedrock devnet launcher')
parser.add_argument('--monorepo-dir', help='Directory of the monorepo', default=os.getcwd())
parser.add_argument('--allocs', help='Only create the allocs and exit', type=bool, action=argparse.BooleanOptionalAction)
parser.add_argument('--test', help='Tests the deployment, must already be deployed', type=bool, action=argparse.BooleanOptionalAction)
parser.add_argument('--l2', help='Which L2 to run', type=str, default='op1')
parser.add_argument('--l2-provider-url', help='URL for the L2 RPC node', type=str, default='http://localhost:19545')
parser.add_argument('--deploy-l2', help='Deploy the L2 onto a running L1 and sequencer network', type=bool, action=argparse.BooleanOptionalAction)
parser.add_argument('--deploy-config', help='Deployment config, relative to packages/contracts-bedrock/deploy-config', default='devnetL1.json')
parser.add_argument('--deploy-config-template', help='Deployment config template, relative to packages/contracts-bedrock/deploy-config', default='devnetL1-nodekit-template.json')
parser.add_argument('--deployment', help='Path to deployment output files, relative to packages/contracts-bedrock/deployments', default='devnetL1')
parser.add_argument('--devnet-dir', help='Output path for devnet config, relative to --monorepo-dir', default='.devnet')
parser.add_argument('--nodekit', help='Run on NodeKit SEQ', type=bool, action=argparse.BooleanOptionalAction)
parser.add_argument("--compose-file", help="Compose file to use for demo images", type=str, default="docker-compose.yml")
parser.add_argument('--eth-pos-dir', help="directory to store repository `eth-pos-devnet`", default='eth-pos-devnet')
parser.add_argument('--jwt-secret', help='jwt secret to access geth http api', type=str, default='0xfad2709d0bb03bf0e8ba3c99bea194575d3e98863133d1af638ed056d1d59345')
parser.add_argument('--zk-dir', help='nodekit-zk directory', type=str, default='nodekit-zk')
parser.add_argument('--l1-rpc-url', help='l1 rpc url', type=str, default='http://localhost:8545')
parser.add_argument('--l1-ws-url', help='l1 ws url', type=str, default='ws://localhost:8546')
parser.add_argument('--launch-l2', help='if launch l2', type=bool, action=argparse.BooleanOptionalAction)
parser.add_argument('--launch-nodekit-l1', help='if launch nodekit l1', type=bool, action=argparse.BooleanOptionalAction)
parser.add_argument('--nodekit-l1-dir', help='directory of nodekit-l1', type=str, default='nodekit-l1')
parser.add_argument('--nodekit-contract', help='nodekit commitment contract address on l1', type=str, default='')
parser.add_argument('--seq-url',  help='seq url', type=str, default='http://127.0.0.1:37029/ext/bc/56iQygPt5wrSCqZSLVwKyT7hAEdraXqDsYqWtWoAWaZSKDSDm')
parser.add_argument('--l1-chain-id', help='chain id of l1', type=str, default='32382')
parser.add_argument('--l2-chain-id', help='chain id of l2', type=str, default='45200')
parser.add_argument('--deploy-contracts', help='deploy contracts for l2 and nodekit-zk', type=bool, action=argparse.BooleanOptionalAction)
parser.add_argument('--mnemonic-words', help='mnemonic words to deploy nodekit-zk contract', type=str, default='test test test test test test test test test test test junk')


# Global environment variables
DEVNET_NO_BUILD = os.getenv('DEVNET_NO_BUILD') == "true"
DEVNET_FPAC = os.getenv('DEVNET_FPAC') == "true"
DEVNET_PLASMA = os.getenv('DEVNET_PLASMA') == "true"

log = logging.getLogger()

class Bunch:
    def __init__(self, **kwds):
        self.__dict__.update(kwds)

class ChildProcess:
    def __init__(self, func, *args):
        self.errq = Queue()
        self.process = Process(target=self._func, args=(func, args))

    def _func(self, func, args):
        try:
            func(*args)
        except Exception as e:
            self.errq.put(str(e))

    def start(self):
        self.process.start()

    def join(self):
        self.process.join()

    def get_error(self):
        return self.errq.get() if not self.errq.empty() else None


def main():
    args = parser.parse_args()

    monorepo_dir = os.path.abspath(args.monorepo_dir)
    devnet_dir = pjoin(monorepo_dir, args.devnet_dir)
    contracts_bedrock_dir = pjoin(monorepo_dir, 'packages', 'contracts-bedrock')
    deployment_dir = pjoin(contracts_bedrock_dir, 'deployments', args.deployment)
    op_node_dir = pjoin(args.monorepo_dir, 'op-node')
    ops_bedrock_dir = pjoin(monorepo_dir, 'ops-bedrock')
    deploy_config_dir = pjoin(contracts_bedrock_dir, 'deploy-config'),
    devnet_config_path = pjoin(contracts_bedrock_dir, 'deploy-config', args.deploy_config)
    devnet_config_template_path = pjoin(contracts_bedrock_dir, 'deploy-config', args.deploy_config_template)
    ops_chain_ops = pjoin(monorepo_dir, 'op-chain-ops')
    sdk_dir = pjoin(monorepo_dir, 'packages', 'sdk')
    eth_pos_dir: str = args.eth_pos_dir
    zk_dir: str = args.zk_dir
    nodekit_l1_dir: str = args.nodekit_l1_dir
    l2_chain_id: int = int(args.l2_chain_id)

    jwt_secret: str = args.jwt_secret

    l1_rpc_url: str = args.l1_rpc_url
    launch_l2: bool = args.launch_l2
    launch_nodekit_l1: bool = args.launch_nodekit_l1
    _deploy_contracts: bool = args.deploy_contracts


    paths = Bunch(
      mono_repo_dir=monorepo_dir,
      devnet_dir=devnet_dir,
      contracts_bedrock_dir=contracts_bedrock_dir,
      deployment_dir=deployment_dir,
      l1_deployments_path=pjoin(deployment_dir, '.deploy'),
      deploy_config_dir=deploy_config_dir,
      devnet_config_path=devnet_config_path,
      devnet_config_template_path=devnet_config_template_path,
      op_node_dir=op_node_dir,
      ops_bedrock_dir=ops_bedrock_dir,
      ops_chain_ops=ops_chain_ops,
      sdk_dir=sdk_dir,
      genesis_l1_path=pjoin(devnet_dir, 'genesis-l1.json'),
      genesis_l2_path=pjoin(devnet_dir, 'genesis-l2.json'),
      allocs_path=pjoin(devnet_dir, 'allocs-l1.json'),
      addresses_json_path=pjoin(devnet_dir, 'addresses.json'),
      sdk_addresses_json_path=pjoin(devnet_dir, 'sdk-addresses.json'),
      rollup_config_path=pjoin(devnet_dir, 'rollup.json'),
      zk_dir=zk_dir,
      l1_rpc_url=l1_rpc_url,
      nodekit_l1_dir=nodekit_l1_dir
    )

    # if args.test:
    #     log.info('Testing deployed devnet')
    #     devnet_test(paths, args.l2_provider_url)
    #     return

    os.makedirs(devnet_dir, exist_ok=True)

    # if args.allocs:
    #     devnet_l1_genesis(paths, args.deploy_config)
    #     return

    # log.info('launching eth pos net')
    # launch_eth_devnet(eth_pos_dir)

    # try:
    #     ethnet_up = wait_up(8545)
    #     if ethnet_up:
    #         log.info("eth net launched successfully")
    # except Exception as e:
    #     log.error(f'unable to launch eth net, waiting failed with error {e}')

    # priv = get_nodekit_zk_contract_deployer_priv()
    # print(priv)
    # return

    # TODO: to be removed
    # if launch_nodekit_l1:
    #     log.info('launching nodekit l1')
    #     deploy_nodekit_i1(paths, args)
    #     return

    if launch_l2:
        log.info('launching op stack')
        devnet_deploy(paths, args)
        return

    if _deploy_contracts:
        try:
            init_devnet_l1_deploy_config(paths, update_timestamp=True, l2_chain_id=l2_chain_id)
            deploy_contracts(paths, args, args.deploy_config, True, jwt_secret)
            log.info('contracts deployed')
        except Exception as e:
            log.error(f'unable to deploy contracts: {e}')

        return


def deploy_nodekit_i1(paths, args):
    nodekit_l1_dir: str = paths.nodekit_l1_dir
    seq_url: str = args.seq_url
    seq_chain_id: str = seq_url.split('/')[-1]
    commitment_contract_addr: str = get_nodekit_zk_contract_addr(paths, args)
    # TODO: generate by phrase
    commitment_contract_wallet: str = 'ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80'
    l1_chain_id: str = args.l1_chain_id
    l1_rpc: str = args.l1_rpc_url

    env = {
        'SEQ_ADDR': seq_url,
        'CHAIN_ID': seq_chain_id,
        'CONTRACT_ADDR': commitment_contract_addr,
        # lstrip or `Failed to convert from hex to ECDSA: invalid hex character 'x' in private key`
        'CONTRACT_WALLET': commitment_contract_wallet.lstrip('0x'),
        'CHAIN_ID_L1': l1_chain_id,
        'L1_RPC': l1_rpc
    }

    log.info(f'using config to deploy nodekit l1: {env}')

    subprocess.Popen([
        'docker', 'compose', 'up', '-d'
    ], cwd=nodekit_l1_dir, env=env)

def launch_eth_devnet(eth_pos_dir: str):
    if not os.path.exists(eth_pos_dir) or not os.path.isdir(eth_pos_dir):
        raise Exception(f"Designated eth_pos_dir not valid: {eth_pos_dir}")

    # run_command([
    #     'sudo', 'bash', 'clean.sh', '-y'
    # ], cwd=eth_pos_dir)

    subprocess.Popen([
        'docker', 'compose', 'up', '-d'
    ], cwd=eth_pos_dir)

def stop_eth_devnet(eth_pos_dir: str):
    run_command([
        'sudo', 'bash', './clean.sh'
    ], cwd=eth_pos_dir)
    run_command([
        'docker', 'compose', 'down'
    ], cwd=eth_pos_dir)

# # TODO: fix dependency for nix-shell
# def get_nodekit_zk_contract_deployer_priv(paths, args) -> str:
#     mnemonic_words = args.mnemonic_words

#     # following from: https://github.com/meherett/python-hdwallet
#     MNEMONIC: str = mnemonic_words
#     # Secret passphrase/password for mnemonic
#     PASSPHRASE: Optional[str] = None  # "meherett"

#     # Initialize Ethereum mainnet BIP44HDWallet
#     bip44_hdwallet: BIP44HDWallet = BIP44HDWallet(cryptocurrency=EthereumMainnet)
#     # Get Ethereum BIP44HDWallet from mnemonic
#     bip44_hdwallet.from_mnemonic(
#         mnemonic=MNEMONIC, language="english", passphrase=PASSPHRASE
#     )
#     # Clean default BIP44 derivation indexes/paths
#     bip44_hdwallet.clean_derivation()

#     print("Mnemonic:", bip44_hdwallet.mnemonic())
#     print("Base HD Path:  m/44'/60'/0'/0/{address_index}", "\n")

#     # Derivation from Ethereum BIP44 derivation path
#     bip44_derivation: BIP44Derivation = BIP44Derivation(
#         cryptocurrency=EthereumMainnet, account=0, change=False, address=0
#     )
#     # Drive Ethereum BIP44HDWallet
#     bip44_hdwallet.from_path(path=bip44_derivation)
#     return bip44_derivation.private_key()


def get_nodekit_zk_contract_addr(paths, args) -> str:
    zk_dir = paths.zk_dir
    l1_chain_id = args.l1_chain_id
    latest_run_path = os.path.join(zk_dir, f'broadcast/Sequencer.s.sol/{l1_chain_id}/run-latest.json')

    with open(latest_run_path, 'r') as f:
        runinfo_str = f.read()
        runinfo = json.loads(runinfo_str)

        return runinfo['transactions'][0]['contractAddress']

# deploy Create2, Nodekit-ZK
def deploy_contracts(paths, args, deploy_config: str, deploy_l2: bool, jwt_secret: str = ''):
    rpc_url = paths.l1_rpc_url
    # rpc_url = 'http://localhost:8545'
    wait_for_rpc_server_local(rpc_url, jwt_secret)
    res = eth_accounts_local(rpc_url, jwt_secret)
    account = res['result'][0]
    log.info(f'Deploying with {account}')
    mnemonic_words = args.mnemonic_words
    sequencer_contract_addr: str = args.nodekit_contract

    # # wait transaction indexing service to be available
    # time.sleep(30)

    # The create2 account is shared by both L2s, so don't redeploy it unless necessary
    # We check to see if the create2 deployer exists by querying its balance

    res = run_command(
        ["cast", "balance", "0x3fAB184622Dc19b6109349B94811493BF2a45362", "--rpc-url", rpc_url],
        capture_output=True,
    )
    deployer_balance = int(res.stdout.strip())
    if deployer_balance == 0:
        # send some ether to the create2 deployer account
        run_command([
            'cast', 'send', '--from', account,
            '--rpc-url', rpc_url,
            '--unlocked', '--value', '1ether', '0x3fAB184622Dc19b6109349B94811493BF2a45362',
            # '--jwt-secret', jwt_secret
        ], env={}, cwd=paths.contracts_bedrock_dir)

        # deploy the create2 deployer
        run_command([
            'cast', 'publish', '--rpc-url', rpc_url,
            '0xf8a58085174876e800830186a08080b853604580600e600039806000f350fe7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe03601600081602082378035828234f58015156039578182fd5b8082525050506014600cf31ba02222222222222222222222222222222222222222222222222222222222222222a02222222222222222222222222222222222222222222222222222222222222222',
            # '--jwt-secret', jwt_secret
        ], env={}, cwd=paths.contracts_bedrock_dir)


    deploy_env = {
        'DEPLOYMENT_CONTEXT': deploy_config.removesuffix('.json'),
        'MNEMONIC': mnemonic_words
    }
    print(f'deploy_env: {deploy_env}')
    if deploy_l2:
        # If deploying an L2 onto an existing L1, use a different deployer salt so the contracts
        # will not collide with those of the existing L2.
        deploy_env['IMPL_SALT'] = os.urandom(32).hex()

    fqn = 'scripts/Deploy.s.sol:Deploy'
    run_command([
        'forge', 'script', fqn, '--sender', account,
        '--rpc-url', rpc_url, '--broadcast',
        '--unlocked'
    ], env=deploy_env, cwd=paths.contracts_bedrock_dir)

    # TODO: to be removed since sequencerContractAddr is not used by either op-node or other op contracts
    # # deploy nodekit-zk contracts
    # #TODO add back later
    # run_command([
    #     'forge', 'script', 'DeploySequencer', '--broadcast',
    #     '--rpc-url', rpc_url,
    # ], env=deploy_env, cwd=paths.zk_dir)

    # update config for l2
    # or will lead to unable to verify l2 blocks
    # sequencer_contract_addr = get_nodekit_zk_contract_addr(paths, args)
    devnetL1_conf = read_json(paths.devnet_config_path)
    devnetL1_conf['nodekitContractAddress'] = sequencer_contract_addr
    write_json(paths.devnet_config_path, devnetL1_conf)

    shutil.copy(paths.l1_deployments_path, paths.addresses_json_path)

    # TODO: to be removed, seems in the recent version of optimisim, they removed sync from Deployer contract
    # log.info('Syncing contracts.')
    # run_command([
    #     'forge', 'script', fqn, '--sig', 'sync()',
    #     '--rpc-url', rpc_url
    # ], env=deploy_env, cwd=paths.contracts_bedrock_dir)

def init_devnet_l1_deploy_config(paths, update_timestamp=False, l2_chain_id=45200):
    deploy_config = read_json(paths.devnet_config_template_path)
    if update_timestamp:
        deploy_config['l1GenesisBlockTimestamp'] = '{:#x}'.format(int(time.time()))
    if DEVNET_FPAC:
        deploy_config['useFaultProofs'] = True
        deploy_config['faultGameMaxDuration'] = 10
    if DEVNET_PLASMA:
        deploy_config['usePlasma'] = True

    deploy_config['l2ChainID'] = l2_chain_id
    write_json(paths.devnet_config_path, deploy_config)

# unused
def devnet_l1_genesis(paths, deploy_config: str):
    log.info('Generating L1 genesis state')

    # Abort if there is an existing geth process listening on localhost:8545. It
    # may cause the op-node to fail to start due to a bad genesis block.
    geth_up = False
    try:
        geth_up = wait_up(8545, retries=1, wait_secs=0)
    except:
        pass
    if geth_up:
        raise Exception('Existing process is listening on localhost:8545, please kill it and try again. (e.g. `pkill geth`)')

    init_devnet_l1_deploy_config(paths)

    geth = subprocess.Popen([
        'geth', '--dev', '--http', '--http.api', 'eth,debug',
        '--verbosity', '4', '--gcmode', 'archive', '--dev.gaslimit', '30000000',
        '--rpc.allow-unprotected-txs'
    ])

    try:
        forge = ChildProcess(deploy_contracts, paths, deploy_config, False)
        forge.start()
        forge.join()
        err = forge.get_error()
        if err:
            raise Exception(f"Exception occurred in child process: {err}")

        res = debug_dumpBlock('devnet.nodekit.xyz')
        response = json.loads(res)
        allocs = response['result']

        write_json(paths.allocs_path, allocs)
    finally:
        geth.terminate()


# Bring up the devnet where the contracts are deployed to L1
def devnet_deploy(paths, args):
    nodekit = args.nodekit
    l2 = args.l2
    l2_chain_id = int(args.l2_chain_id)
    # which will be prepended to names of docker volumnes and services so we can run several rollups
    composer_project_name = f'op-devnet_{l2_chain_id}'
    l2_provider_url = args.l2_provider_url
    compose_file = args.compose_file
    l1_rpc_url = args.l1_rpc_url
    l1_ws_url = args.l1_ws_url
    seq_addr: str = args.seq_url
    seq_chain_id = seq_addr.split('/')[-1]

    conf = {
        l2_provider_url,
        compose_file,
        l1_rpc_url,
        l1_ws_url,
        seq_addr,
        seq_chain_id
    }

    print(f'using config {conf}')

    # TODO: to be removed since we don't need to launch l2 ourselves
    # if os.path.exists(paths.genesis_l1_path) and os.path.isfile(paths.genesis_l1_path):
    #     log.info('L1 genesis already generated.')
    # elif not args.deploy_l2:
    #     # Generate the L1 genesis, unless we are deploying an L2 onto an existing L1.
    #     log.info('Generating L1 genesis.')
    #     if os.path.exists(paths.allocs_path) == False:
    #         devnet_l1_genesis(paths, args.deploy_config)

    #     # It's odd that we want to regenerate the devnetL1.json file with
    #     # an updated timestamp different than the one used in the devnet_l1_genesis
    #     # function.  But, without it, CI flakes on this test rather consistently.
    #     # If someone reads this comment and understands why this is being done, please
    #     # update this comment to explain.
    #     init_devnet_l1_deploy_config(paths, update_timestamp=True)
    #     outfile_l1 = pjoin(paths.devnet_dir, 'genesis-l1.json')
    #     run_command([
    #         'go', 'run', 'cmd/main.go', 'genesis', 'l1',
    #         '--deploy-config', paths.devnet_config_path,
    #         '--l1-allocs', paths.allocs_path,
    #         '--l1-deployments', paths.addresses_json_path,
    #         '--outfile.l1', outfile_l1,
    #     ], cwd=paths.op_node_dir)

    # if args.deploy_l2:
    #     # L1 and sequencer already exist, just create the deploy config and deploy the L1 contracts
    #     # for the new L2.
    #     init_devnet_l1_deploy_config(paths, update_timestamp=True)
    #     deploy_contracts(paths, args.deploy_config, args.deploy_l2)
    # else:
    #     # Deploy L1 and sequencer network.
    #     log.info('Starting L1.')
    #     run_command(['docker', 'compose', '-f', compose_file, 'up', '-d', 'l1'], cwd=paths.ops_bedrock_dir, env={
    #         'PWD': paths.ops_bedrock_dir,
    #         'DEVNET_DIR': paths.devnet_dir
    #     })
    #     #wait_up(8545)
    #     wait_for_rpc_server('devnet.nodekit.xyz')

    #     log.info('Bringing up `artifact-server`')
    #     run_command(['docker', 'compose', 'up', '-d', 'artifact-server'], cwd=paths.ops_bedrock_dir, env={
    #         'PWD': paths.ops_bedrock_dir,
    #         'DEVNET_DIR': paths.devnet_dir
    #     })



    # Re-build the L2 genesis unconditionally in NodeKit mode, since we require the timestamps to be recent.
    # if not nodekit and os.path.exists(paths.genesis_l2_path) and os.path.isfile(paths.genesis_l2_path):
    if os.path.exists(paths.genesis_l2_path) and os.path.isfile(paths.genesis_l2_path):
        log.info('L2 genesis and rollup configs already generated.')
    else:
        log.info('Generating L2 genesis and rollup configs.')
        run_command([
            'go', 'run', 'cmd/main.go', 'genesis', 'l2',
            '--l1-rpc', l1_rpc_url,
            '--deploy-config', paths.devnet_config_path,
            '--l1-deployments', pjoin(paths.deployment_dir, '.deploy'),
            '--outfile.l2', pjoin(paths.devnet_dir, 'genesis-l2.json'),
            '--outfile.rollup', pjoin(paths.devnet_dir, 'rollup.json')
        ], cwd=paths.op_node_dir)

    rollup_config = read_json(paths.rollup_config_path)
    addresses = read_json(paths.addresses_json_path)
    l2_provider_port = int(l2_provider_url.split(':')[-1])
    l2_provider_http = l2_provider_url

    log.info(f'l2 provider http: {l2_provider_http}, port: {l2_provider_port}')

    log.info('Bringing up L2.')
    run_command(['docker', 'compose', '-f', compose_file, 'up', '-d', f'{l2}-l2', f'{l2}-geth-proxy'], cwd=paths.ops_bedrock_dir, env={
        'PWD': paths.ops_bedrock_dir,
        'DEVNET_DIR': paths.devnet_dir,
        'SEQ_ADDR': seq_addr,
        'SEQ_CHAIN_ID': seq_chain_id,
        'OP1_L2_RPC_PORT': str(l2_provider_port),
        'COMPOSE_PROJECT_NAME': composer_project_name
    })

    # l2_provider_port = int(l2_provider_url.split(':')[-1])
    # l2_provider_http = l2_provider_url.removeprefix('http://')
    wait_up(l2_provider_port)
    wait_for_rpc_server_local(l2_provider_http)

    l2_output_oracle = addresses['L2OutputOracleProxy']
    log.info(f'Using L2OutputOracle {l2_output_oracle}')
    batch_inbox_address = rollup_config['batch_inbox_address']
    log.info(f'Using batch inbox {batch_inbox_address}')

    log.info('Bringing up `op-node`, `op-proposer` and `op-batcher`.')
    command = ['docker', 'compose', '-f', compose_file, 'up', '-d']
    if args.deploy_l2:
        # If we are deploying onto an existing L1, don't restart the services that are already
        # running.
        command.append('--no-recreate')
    services = [f'{l2}-node', f'{l2}-proposer', f'{l2}-batcher']
    run_command(command + services, cwd=paths.ops_bedrock_dir, env={
        'PWD': paths.ops_bedrock_dir,
        'L2OO_ADDRESS': l2_output_oracle,
        'SEQUENCER_BATCH_INBOX_ADDRESS': batch_inbox_address,
        'DEVNET_DIR': paths.devnet_dir,
        'SEQ_ADDR': seq_addr,
        'SEQ_CHAIN_ID': seq_chain_id,
        'L1WS': l1_ws_url,
        'L1RPC': l1_rpc_url,
        'COMPOSE_PROJECT_NAME': composer_project_name
    })

    #log.info('Starting block explorer')
    #run_command(['docker-compose', '-f', compose_file, 'up', '-d', f'{l2}-blockscout'], cwd=paths.ops_bedrock_dir)

    log.info('Devnet ready.')

def eth_accounts_local(url, jwt_secret=None):
    log.info(f'Fetch eth_accounts {url}')
    headers = {
        "Content-Type": "application/json"
    }
    payload = {"id":2, "jsonrpc":"2.0", "method": "eth_accounts", "params":[]}
    if jwt_secret:
        token = generate_jwt_token(jwt_secret)
        headers['Authorization'] = f"Bearer {token}"

    session = requests.Session()
    resp = session.post(url, json=payload, headers=headers)
    session.close()

    return resp.json()

def eth_accounts(url):
    log.info(f'Fetch eth_accounts {url}')
    conn = http.client.HTTPSConnection(url)
    headers = {'Content-type': 'application/json'}
    body = '{"id":2, "jsonrpc":"2.0", "method": "eth_accounts", "params":[]}'
    conn.request('POST', '/', body, headers)
    response = conn.getresponse()
    data = response.read().decode()
    conn.close()
    return data


def debug_dumpBlock_local(url, jwt_secret=None):
    log.info(f'Fetch debug_dumpBlock {url}')
    conn = http.client.HTTPConnection(url)
    headers = {'Content-type': 'application/json'}
    body = '{"id":3, "jsonrpc":"2.0", "method": "debug_dumpBlock", "params":["latest"]}'
    conn.request('POST', '/', body, headers)
    response = conn.getresponse()
    data = response.read().decode()
    conn.close()
    return data
    # headers = {
    #     "Content-Type": "application/json"
    # }
    # payload = {"id":3, "jsonrpc":"2.0", "method": "debug_dumpBlock", "params":["latest"]}
    # if jwt_secret:
    #     token = generate_jwt_token(jwt_secret)
    #     headers['Authorization'] = f"Bearer {token}"

    # session = requests.Session()
    # resp = session.post(url, json=payload, headers=headers)
    # session.close()
    # return resp.json()

def debug_dumpBlock(url):
    log.info(f'Fetch debug_dumpBlock {url}')
    conn = http.client.HTTPSConnection(url)
    headers = {'Content-type': 'application/json'}
    body = '{"id":3, "jsonrpc":"2.0", "method": "debug_dumpBlock", "params":["latest"]}'
    conn.request('POST', '/', body, headers)
    response = conn.getresponse()
    data = response.read().decode()
    conn.close()
    return data


def wait_for_rpc_server(url):
    log.info(f'Waiting for RPC server at {url}')

    conn = http.client.HTTPSConnection(url)
    headers = {'Content-type': 'application/json'}
    body = '{"id":1, "jsonrpc":"2.0", "method": "eth_chainId", "params":[]}'

    while True:
        try:
            conn.request('POST', '/', body, headers)
            response = conn.getresponse()
            conn.close()
            log.info(response)
            if response.status < 300:
                log.info(f'RPC server at {url} ready')
                return
        except Exception as e:
            log.info(f'Error connecting to RPC: {e}')
            log.info(f'Waiting for RPC server at {url}')
            time.sleep(1)


def wait_for_rpc_server_local(url, jwt_secret=None):
    headers = {
        "Content-Type": "application/json"
    }
    payload = {"id":1, "jsonrpc":"2.0", "method": "eth_chainId", "params":[]}
    if jwt_secret:
        token = generate_jwt_token(jwt_secret)
        headers['Authorization'] = f"Bearer {token}"

    while True:
        session = requests.Session()
        try:
            resp = session.post(url, json=payload, headers=headers)
            if resp.status_code < 300:
                log.info(f'RPC server at {url} ready')
                session.close()
                return
            session.close()
        except Exception as e:
            time.sleep(1)
            log.warn(f'unable to connect to geth {e}')

def generate_jwt_token(secret, expiration_seconds=JWT_EXPIRATION_SECONDS):
    payload = {
        "exp": int(time.time()) + expiration_seconds
    }
    token = jwt.encode(payload, secret, algorithm='HS256')
    return token

def deploy_erc20(paths, l2_provider_url):
    run_command(
         ['npx', 'hardhat',  'deploy-erc20', '--network',  'devnetL1', '--l2-provider-url', l2_provider_url],
         cwd=paths.sdk_dir,
         timeout=60,
    )

CommandPreset = namedtuple('Command', ['name', 'args', 'cwd', 'timeout'])

def devnet_test(paths, l2_provider_url):
    # Check the L2 config
    run_command(
        ['go', 'run', 'cmd/check-l2/main.go', '--l2-rpc-url', l2_provider_url, '--l1-rpc-url', 'https://devnet.nodekit.xyz'],
        cwd=paths.ops_chain_ops,
    )

    # Run the commands with different signers, so the ethereum nonce management does not conflict
    # And do not use devnet system addresses, to avoid breaking fee-estimation or nonce values.
    run_commands([
        CommandPreset('erc20-test',
          ['npx', 'hardhat',  'deposit-erc20', '--network',  'devnetL1',
           '--l1-contracts-json-path', paths.addresses_json_path, '--l2-provider-url', l2_provider_url, '--signer-index', '14'],
          cwd=paths.sdk_dir, timeout=8*60),
        CommandPreset('eth-test',
          ['npx', 'hardhat',  'deposit-eth', '--network',  'devnetL1',
           '--l1-contracts-json-path', paths.addresses_json_path, '--l2-provider-url', l2_provider_url, '--signer-index', '15'],
          cwd=paths.sdk_dir, timeout=8*60),
    ], max_workers=2)


def run_commands(commands: list[CommandPreset], max_workers=2):
    with concurrent.futures.ThreadPoolExecutor(max_workers=max_workers) as executor:
        futures = [executor.submit(run_command_preset, cmd) for cmd in commands]

        for future in concurrent.futures.as_completed(futures):
            result = future.result()
            if result:
                print(result.stdout)


def run_command_preset(command: CommandPreset):
    with subprocess.Popen(command.args, cwd=command.cwd,
                          stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True) as proc:
        try:
            # Live output processing
            for line in proc.stdout:
                # Annotate and print the line with timestamp and command name
                timestamp = datetime.datetime.utcnow().strftime('%H:%M:%S.%f')
                # Annotate and print the line with the timestamp
                print(f"[{timestamp}][{command.name}] {line}", end='')

            stdout, stderr = proc.communicate(timeout=command.timeout)

            if proc.returncode != 0:
                raise RuntimeError(f"Command '{' '.join(command.args)}' failed with return code {proc.returncode}: {stderr}")

        except subprocess.TimeoutExpired:
            raise RuntimeError(f"Command '{' '.join(command.args)}' timed out!")

        except Exception as e:
            raise RuntimeError(f"Error executing '{' '.join(command.args)}': {e}")

        finally:
            # Ensure process is terminated
            proc.kill()
    return proc.returncode

def run_command(args, check=True, shell=False, cwd=None, env=None, timeout=None, capture_output=False):
    env = env if env else {}
    return subprocess.run(
        args,
        check=check,
        shell=shell,
        capture_output=capture_output,
        env={
            **os.environ,
            **env
        },
        cwd=cwd,
        timeout=timeout
    )


def wait_up(port, retries=10, wait_secs=1):
    for i in range(0, retries):
        log.info(f'Trying 127.0.0.1:{port}')
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        try:
            s.connect(('127.0.0.1', int(port)))
            s.shutdown(2)
            log.info(f'Connected 127.0.0.1:{port}')
            return True
        except Exception:
            time.sleep(wait_secs)

    raise Exception(f'Timed out waiting for port {port}.')


def write_json(path, data):
    with open(path, 'w+') as f:
        json.dump(data, f, indent='  ')


def read_json(path):
    with open(path, 'r') as f:
        return json.load(f)
