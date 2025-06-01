import logging
import requests
import time

from bs4 import BeautifulSoup as bs
from typing import List

logger = logging.getLogger("main")


def get_nations(user: str, region: str, use_cache: bool, rate: int) -> List[str]:
    headers = {"User-Agent": user}

    url = f"https://www.nationstates.net/cgi-bin/api.cgi?region={region}&q=nations"

    nations = bs(requests.get(url, headers=headers).text, "xml").find("NATIONS")

    if not nations:
        raise ValueError("nations not found")

    nations = nations.text.split(":")

    logger.info("total number of nations: %d", len(nations))

    try:
        logger.debug("checking exclusions.txt")
        with open("exclusions.txt") as ifile:
            exclusions = ifile.read().splitlines()

            if use_cache:
                nations = [nation for nation in nations if nation not in exclusions]
    except FileNotFoundError:
        logger.debug("exclusions.txt not found, creating")
        open("./exclusions.txt", "w").close()

    logger.info("nations to check after exclusions: %d", len(nations))

    result = []

    logger.info(
        "checking recruitment permissions for each nation, this may take a while"
    )

    for nation in nations:
        logger.debug("checking nation: %s", nation)
        try:
            if recruitment_enabled(user, nation):
                result.append(nation)
        except Exception:
            logger.exception("unable to retrieve data for nation: %s", nation)
        time.sleep(30 / rate)

    logger.info("nations with recruitment enabled: %d", len(result))

    with open("exclusions.txt", "a") as ofile:
        ofile.write("\n".join([nation for nation in nations if nation not in result]))
        ofile.write("\n")

    return result


def recruitment_enabled(user: str, nation: str) -> bool:
    headers = {"User-Agent": user}

    url = f"https://www.nationstates.net/cgi-bin/api.cgi?nation={nation}&q=tgcanrecruit"

    result = bs(requests.get(url, headers=headers).text, "xml").find("TGCANRECRUIT")

    if not result:
        raise ValueError("tgcanrecruit not found")

    return result.text == "1"
