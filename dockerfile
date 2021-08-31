FROM docker.io/python:3

WORKDIR /webcache

ADD ./requirements.txt .
RUN pip3 install -r requirements.txt

ADD . ./

CMD ["flask", "run"]