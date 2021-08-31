import os
import logging
import threading
import constants
import config
import api

from dotenv import load_dotenv

if __name__ == "__main__":
    load_dotenv()
    logging.basicConfig(format=constants.LOG_FORMAT, level=logging.INFO,
                        datefmt=constants.TIME_FORMAT)

    config_path = os.getenv(constants.ENV_CONFIG_PATH)
    sleep_time = os.getenv(constants.ENV_SLEEP_TIME)
    if sleep_time:
        sleep_time = float(sleep_time)
    else:
        sleep_time = 0.

    fn_args = (config_path, sleep_time)
    fn_config = threading.Thread(target=config.run, args=fn_args, daemon=True)
    fn_config.start()

    server_port = os.getenv(constants.ENV_SERVER_PORT)
    api.start(server_port)