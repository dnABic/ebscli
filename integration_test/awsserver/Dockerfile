FROM dnabic/moto:v0.2.0
MAINTAINER dnABic

#COPY test_env/aws-moto-server/requirements.txt /opt/venv/requirements.txt
#RUN /opt/venv/bin/pip install -r /opt/venv/requirements.txt

COPY .boto /root/.boto

#CMD cd /opt/venv/ && /opt/venv/bin/moto_server -p3000

ENV AWS_ACCESS_KEY_ID the_key
ENV AWS_SECRET_ACCESS_KEY the_secret

#EXPOSE 3000
CMD["moto_server", "ec2", "-H", "0.0.0.0"]
