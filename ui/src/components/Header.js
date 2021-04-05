import React from 'react'
import { fetchDecks } from '../services/deck.js'
import { cx, css } from 'pretty-lights'
import PropTypes from 'prop-types'
import { B, U, R, G, W } from '../Mana.js'
import { MTGFailLogoWhite as MTGFailLogo } from '../MTGFailLogo.js'
import { Link } from 'react-router-dom'

const box = css`
  display: flex;
  justify-content: space-evenly;
  background-color: black;
`
const style = css`
  flex: 1 1 auto;
  display: flex;
  justify-content: space-evenly;
  align-items: center;
  background-color: black;
`
const logo = css`
  margin-left: 1em;
  height: 100px;
  width: 100px;
`
const image = css`
  height: 45px;
  transition: transform 0.15s;
  width: 45px;
  z-offset: 99;
`
const grow = css`
  &:hover {
    transform: scale(1.1);
    border-radius: 100%;
  }
`
const log = css`
  margin-left: auto;
`

const OMNOM = 'https://deckbox.org/sets/2785835'
const SHORT_FAIRY = 'https://deckbox.org/sets/2811132'

const Header = ({ setDeck, setTTSDeck, login, onError }) => {
  const [reload, setReload] = React.useState(false)

  React.useEffect(() => {
    if (!reload) return
    const f = async () => {
      try {
        const decks = await fetchDecks(reload, onError)

        setTTSDeck(decks.tts)
        setDeck(decks.internal.sort((a, b) => (a.name > b.name ? 1 : -1)))
      } catch (e) {
        onError(e)
      }
    }
    f()
    setReload(false)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [reload])

  return (
    <div className={box}>
      <Link to="/">
        <MTGFailLogo style={logo} />
      </Link>
      <div className={style}>
        <span onClick={() => setReload(SHORT_FAIRY)}>
          <B style={cx(grow, image)} />
        </span>
        <U style={image} />
        <W style={image} />
        <span onClick={() => setReload(OMNOM)}>
          <R style={cx(grow, image)} />
        </span>
        <G style={image} />
        {login && <div className={log}>{login} </div>}
      </div>
    </div>
  )
}

Header.propTypes = {
  login: PropTypes.element,
}
export default Header
