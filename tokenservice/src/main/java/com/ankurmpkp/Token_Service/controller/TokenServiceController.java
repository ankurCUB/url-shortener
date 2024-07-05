package com.ankurmpkp.Token_Service.controller;

import com.ankurmpkp.Token_Service.dto.TokensDto;
import com.ankurmpkp.Token_Service.service.ITokenService;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping(path = "api", produces = {MediaType.APPLICATION_JSON_VALUE})
public class TokenServiceController {

    private final ITokenService iTokenService;

    public TokenServiceController(ITokenService iTokenService){
        this.iTokenService = iTokenService;
    }

    @PostMapping("/create")
    ResponseEntity<TokensDto> createToke(){
        TokensDto tokensDto = iTokenService.createToken();
        return ResponseEntity.status(HttpStatus.CREATED).body(tokensDto);
    }

}
