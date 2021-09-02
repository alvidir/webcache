import os
import time
import logging

def run(path: str, sleep: float):
    if not path:
        raise Exception("a config directory must be set")

    if not sleep or sleep < 0:
        raise Exception("a positive sleep time must be set")
        
    iteration_count = 0

    while True:
        iteration_count += 1
        logging.info("running config iteration #%s" % iteration_count)

        time.sleep(sleep)