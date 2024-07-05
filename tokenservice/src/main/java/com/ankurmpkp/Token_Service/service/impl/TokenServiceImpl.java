package com.ankurmpkp.Token_Service.service.impl;

import com.ankurmpkp.Token_Service.dto.TokensDto;
import com.ankurmpkp.Token_Service.entity.Tokens;
import com.ankurmpkp.Token_Service.mapper.TokensMapper;
import com.ankurmpkp.Token_Service.repository.TokenRepository;
import com.ankurmpkp.Token_Service.service.ITokenService;
import lombok.AllArgsConstructor;
import org.springframework.stereotype.Service;

@Service
@AllArgsConstructor
public class TokenServiceImpl implements ITokenService {

    private TokenRepository tokenRepository;

    @Override
    public TokensDto createToken() {
        Tokens tokens = new Tokens();
        tokens.isAssigned = true;
        Tokens savedToken = tokenRepository.save(tokens);
        tokenRepository.flush();
        return TokensMapper.mapToTokensDto(savedToken, new TokensDto());
    }
}
