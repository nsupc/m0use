import logging
import requests

from typing import List

logger = logging.getLogger("main")


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

    token = requests.post(url=login_url, json=data).json()["token"]

    headers = {"Authorization": f"Bearer: {token}"}

    telegrams = []

    for nation in nations:
        telegrams.append(
            {
                "sender": "the_europeian_government",
                "id": f"{telegram_id}",
                "secret_key": telegram_key,
                "recipient": nation,
                "tg_type": "standard",
            }
        )

    resp = requests.post(telegram_url, headers=headers, json=telegrams)

    logger.info("telegram request sent: %d", resp.status_code)
