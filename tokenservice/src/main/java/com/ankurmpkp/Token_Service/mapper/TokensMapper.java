package com.ankurmpkp.Token_Service.mapper;

import com.ankurmpkp.Token_Service.dto.TokensDto;
import com.ankurmpkp.Token_Service.entity.Tokens;

public class TokensMapper {
    public static TokensDto mapToTokensDto(Tokens tokens, TokensDto tokensDto){
        tokensDto.setTokenId(tokens.tokenId);
        return tokensDto;
    }
}
