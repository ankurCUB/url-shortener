package com.ankurmpkp.Token_Service.repository;

import com.ankurmpkp.Token_Service.entity.Tokens;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface TokenRepository extends JpaRepository<Tokens, Integer> {

}
