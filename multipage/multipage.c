
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int main(int argc,char **argv)
{
	FILE *fp;
	char *cmd = (char *)malloc(65536);
	char *buffer = (char *)malloc(1024);
	if(argc < 2) {
		printf("args error !\n");
		exit(1);
	}

	strncat(cmd,"google-chrome",strlen("google-chrome"));

	//cmd[strlen(cmd)] = '\0';
	printf("cmd = %s\n",cmd);
	if(NULL == (fp = fopen(argv[1],"r"))) {
		printf("error occur while fopen !\n");
		exit(1);
	}

	while(NULL != fgets(buffer,1024,fp)) {
		//buffer[strlen(buffer) - 1] = '\0';
		strncat(cmd," ",strlen(" "));
		//cmd[strlen(cmd)] = '-';
		strncat(cmd,buffer,strlen(buffer));
		//printf("cmd = %s\n",cmd);
		//printf("buffer : %s\n",buffer);
	}

	//strncat(cmd," ",strlen(" "));
	//strncat(cmd,"&",strlen("&"));
	printf("system cmd : %s\n",cmd);
	system(cmd);
	fclose(fp);
}

