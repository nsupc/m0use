import argparse
import logging
import yaml

from logtail import LogtailHandler
from typing import Literal, Self, Optional


logger = logging.getLogger("main")


class LogConfig:
    token: Optional[str]
    endpoint: Optional[str]
    level: Literal["DEBUG", "INFO", "WARNING", "ERROR"]

    def __init__(
        self,
        token: str | None = None,
        endpoint: str | None = None,
        level: Literal["DEBUG", "INFO", "WARNING", "ERROR"] = "INFO",
    ) -> None:
        if not (token and endpoint):
            self.token = None
            self.endpoint = None

        self.token = token
        self.endpoint = endpoint

        if level in ["DEBUG", "INFO", "WARNING", "ERROR"]:
            self.level = level
        else:
            self.level = "INFO"

        if self.token and self.endpoint:
            handler = LogtailHandler(
                source_token=self.token,
                host=self.endpoint,
            )
        else:
            handler = logging.StreamHandler()

        logger.setLevel(self.level)
        logger.addHandler(handler)


class Eurocore:
    url: str
    user: str
    password: str

    def __init__(self, url: str, user: str, password: str) -> None:
        self.url = url
        self.user = user
        self.password = password

        self.url = self.url.strip("/")


class Telegram:
    id: int
    key: str

    def __init__(self, id: int, key: str) -> None:
        self.id = id
        self.key = key


class Config:
    user: str
    region: str
    telegram: Telegram
    eurocore: Eurocore
    request_rate: int
    cache_results: bool
    log_config: LogConfig

    def __init__(
        self,
        user: str,
        region: str,
        telegram_id: int,
        telegram_key: str,
        eurocore_url: str,
        eurocore_user: str,
        eurocore_password: str,
        request_rate: int,
        token: str | None,
        endpoint: str | None,
        log_level: Literal["DEBUG", "INFO", "WARNING", "ERROR"] = "INFO",
        cache_results: bool = True,
    ):
        telegram = Telegram(telegram_id, telegram_key)
        eurocore = Eurocore(eurocore_url, eurocore_user, eurocore_password)
        log_config = LogConfig(token, endpoint, log_level)

        self.user = user
        self.region = region
        self.telegram = telegram
        self.eurocore = eurocore
        self.request_rate = request_rate
        self.cache_results = cache_results
        self.log_config = log_config

    @classmethod
    def from_yml(cls, path: str = "./config.yml") -> Self:
        with open(path, "r") as in_file:
            data = yaml.safe_load(in_file)

        user = data.get("user")
        region = data.get("region")

        telegram = data.get("telegram")
        telegram_id = telegram.get("id")
        telegram_key = telegram.get("key")

        eurocore = data.get("eurocore")
        eurocore_url = eurocore.get("url")
        eurocore_user = eurocore.get("user")
        eurocore_password = eurocore.get("password")

        request_rate = data.get("request_rate")
        cache_results = data.get("cache_results")

        log = data.get("log")
        token = log.get("token")
        endpoint = log.get("endpoint")
        level = log.get("level")

        return cls(
            user,
            region,
            telegram_id,
            telegram_key,
            eurocore_url,
            eurocore_user,
            eurocore_password,
            request_rate,
            token,
            endpoint,
            level,
            cache_results,
        )
