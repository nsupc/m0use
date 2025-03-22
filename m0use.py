import argparse
import logging
import logging.config
import requests
import time

from bs4 import BeautifulSoup as bs
from typing import List, Literal


class Cli:
    user: str
    region: str
    telegram_id: int
    telegram_key: str
    eurocore_url: str
    eurocore_user: str
    eurocore_password: str
    request_rate: int
    cache_results: bool
    log_level: Literal["DEBUG", "INFO", "WARNING", "ERROR"]


def parse_args() -> Cli:
    parser = argparse.ArgumentParser(
        description="m0use: a eurocore helper script to identify and telegram nations in a region that have recruitment telegrams enabled"
    )

    parser.add_argument("-u", "--user", type=str, required=True, help="User-Agent")

    parser.add_argument(
        "-r", "--region", type=str, required=True, help="NationStates region to check"
    )

    parser.add_argument(
        "-i",
        "--telegram-id",
        type=int,
        required=True,
        help="telegram id to send to nations with recruitment telegrams enabled",
    )

    parser.add_argument(
        "-t",
        "--telegram-key",
        type=str,
        required=True,
        help="telegram secret to send to nations with recruitment enabled",
    )

    parser.add_argument(
        "-e",
        "--eurocore-url",
        type=str,
        required=True,
        help="base url for eurocore instance",
    )

    parser.add_argument(
        "-m", "--eurocore-user", type=str, required=True, help="eurocore username"
    )

    parser.add_argument(
        "-p", "--password", type=str, required=True, help="eurocore password"
    )

    parser.add_argument(
        "-q",
        "--request-rate",
        type=int,
        required=False,
        default=30,
        help="the number of requests to make to NationStates in 30 seconds. Value between 1 and 45.",
    )

    parser.add_argument(
        "-c",
        "--cache-results",
        type=bool,
        required=False,
        default=False,
        choices=[True, False],
        help="output nations with recruitment telegrams disabled to ./exclusions.txt",
    )

    parser.add_argument(
        "-l",
        "--log-level",
        type=str,
        required=False,
        choices=["DEBUG", "INFO", "WARNING", "ERROR"],
        default="INFO",
    )

    cli = parser.parse_args(namespace=Cli)

    cli.region = cli.region.lower().replace(" ", "_")
    cli.eurocore_url = cli.eurocore_url.strip("/")
    if cli.request_rate < 1 or cli.request_rate > 45:
        raise ValueError("request rate must be between 1 and 45")

    return cli


def get_nations(user: str, region: str, use_cache: bool, rate: int) -> List[str]:
    headers = {"User-Agent": user}

    url = f"https://www.nationstates.net/cgi-bin/api.cgi?region={region}&q=nations"

    nations = (
        bs(requests.get(url, headers=headers).text, "xml")
        .find("NATIONS")
        .text.split(":")
    )

    logging.info("total number of nations: %d", len(nations))

    try:
        logging.debug("checking exclusions.txt")
        with open("exclusions.txt") as ifile:
            exclusions = ifile.read().splitlines()

            if use_cache:
                nations = [nation for nation in nations if nation not in exclusions]
    except FileNotFoundError:
        logging.debug("exclusions.txt not found, creating")
        open("./exclusions.txt", "w").close()

    logging.info("nations to check after exclusions: %d", len(nations))

    result = []

    logging.info(
        "checking recruitment permissions for each nation, this may take a while"
    )

    for nation in nations:
        logging.debug("checking nation: %s", nation)
        if recruitment_enabled(user, nation):
            result.append(nation)
        time.sleep(30 / rate)

    logging.info("nations with recruitment enabled: %d", len(result))

    with open("exclusions.txt", "a") as ofile:
        ofile.write("\n".join([nation for nation in nations if nation not in result]))
        ofile.write("\n")

    return result


def recruitment_enabled(user: str, nation: str) -> bool:
    headers = {"User-Agent": user}

    url = f"https://www.nationstates.net/cgi-bin/api.cgi?nation={nation}&q=tgcanrecruit"

    try:
        result = (
            bs(requests.get(url, headers=headers).text, "xml").find("TGCANRECRUIT").text
        )
    except Exception as e:
        logging.exception("error retrieving status for %s", nation, exc_info=e)

    if result == "1":
        return True
    else:
        return False


def send_telegrams(
    nations: List[str],
    eurocore_url: str,
    user: str,
    password: str,
    telegram_id: int,
    telegram_key: str,
):
    login_url = f"{eurocore_url}/login"
    telegram_url = f"{eurocore_url}/telegrams"

    data = {"username": user, "password": password}

    token = requests.post(url=login_url, data=data).json()["token"]

    headers = {"Authentication", f"Bearer: {token}"}

    telegrams = []

    for nation in nations:
        telegrams.append(
            {
                "sender": "the_europeian_government",
                "id": telegram_id,
                "secret_key": telegram_key,
                "recipient": nation,
                "tg_type": "standard",
            }
        )

    resp = requests.post(telegram_url, headers=headers, json=telegrams)

    logging.info("telegram request sent: %d", resp.status_code)


def main():
    cli = parse_args()
    logging.config.fileConfig("logging.conf")
    logging.getLogger().setLevel(cli.log_level)

    attrs = vars(cli)
    logging.info(", ".join("%s: %s" % item for item in attrs.items()))

    nations = get_nations(cli.user, cli.region, cli.cache_results, cli.request_rate)

    send_telegrams(
        nations,
        cli.eurocore_url,
        cli.eurocore_user,
        cli.eurocore_password,
        cli.telegram_id,
        cli.telegram_key,
    )


if __name__ == "__main__":
    main()
