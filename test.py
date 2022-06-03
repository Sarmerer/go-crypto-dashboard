import logging
import time


if __name__ == '__main__':
    while True:
        print('current time: %s' % time.ctime())
        time.sleep(5)
