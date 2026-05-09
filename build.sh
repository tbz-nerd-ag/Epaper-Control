git clone https://github.com/tbz-nerd-ag/Epaper-Untis.git epaper-untis
cd epaper-untis
docker build -t epaper-untis:v1 .
cd ..

git clone https://github.com/tbz-nerd-ag/Epaper-ImageGen.git epaper-imagegen
cd epaper-imagegen
docker build -t epaper-imagegen:v1 .
cd ..

git clone https://github.com/ingressy/Epaper-Control.git epaper-control
cd epaper-control
docker build -t epaper-control:v1 .
cp docker-compose.yml /home/dockeruser/docker-compose.yml
cd ..

docker compose up 
