package com.ankurmpkp.Token_Service.dto;

import lombok.Data;

@Data
public class TokensDto {
    private int start = -1;
    private int end = -1;

    public void convertToRange(int tokenId){
        start = (tokenId-1)*1000+1;
        end = tokenId*1000;
    }
}
