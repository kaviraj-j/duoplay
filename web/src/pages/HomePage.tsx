import { gameApi, roomApi } from "@/api";
import { useRoom } from "@/contexts/RoomContext";
import type { GameListPayload } from "@/types";
import { useEffect, useState } from "react";
import {
  Card,
  CardContent,
  CardActions,
  Typography,
  Button,
  Grid,
  Container,
} from "@mui/material";
import { useNavigate } from "react-router-dom";

const HomePage = () => {
  const [games, setGames] = useState<GameListPayload[]>([]);
  const { room } = useRoom();
  const navigate = useNavigate();

  useEffect(() => {
    gameApi.getGamesList().then((gamesResponse) => {
      setGames(gamesResponse.data);
    });
  }, []);

  const handleChooseGame = (gameName: string) => {
    console.log("Choose game:", gameName);
    roomApi.chooseGame(room!.id, gameName);
  };

  const handleViewGame = (gameName: string) => {
    navigate(`/game/${gameName}`);
  };

  return (
    <Container sx={{ py: 4 }}>
      <Grid container spacing={3}>
        {games.map((game) => (
          <Grid item xs={12} sm={6} md={4} key={game.name}>
            <Card
              sx={{ height: "100%", display: "flex", flexDirection: "column" }}
            >
              <CardContent sx={{ flexGrow: 1 }}>
                <Typography variant="h5" component="h2" gutterBottom>
                  {game.display_name}
                </Typography>
                {/* {game.description && (
                  <Typography variant="body2" color="text.secondary">
                    {game.description}
                  </Typography>
                )} */}
              </CardContent>
              <CardActions sx={{ justifyContent: "flex-end", p: 2 }}>
                <Button
                  size="small"
                  variant="outlined"
                  onClick={() => handleChooseGame(game.name)}
                  disabled={!room}
                >
                  Choose Game
                </Button>
                {/* <Button
                  size="small"
                  variant="contained"
                  onClick={() => handleViewGame(game.name)}
                  disabled={!room}
                >
                  View Game
                </Button> */}
              </CardActions>
            </Card>
          </Grid>
        ))}
      </Grid>
    </Container>
  );
};

export default HomePage;
