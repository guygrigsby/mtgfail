import React from 'react'
import Card from './Card.js'
import { cx, css } from 'pretty-lights'
import './ImageBox.css'

const card = (z, overlap) => {
  return css`
    height: ${overlap ? '50px' : 'auto'};
    z-index: ${z};
    overflow: visible;
    transition: all 0.15s ease-in-out;
    &:hover {
      z-index: ${z + 20};
      transform: scale(102%);
    }}
  `
}
const ImageBox = ({ overlap = true, deck, classes, setSelected, onError }) => {
  return (
    <div className={cx('image-box', classes)} style={{ overflow: 'visible' }}>
      {deck ? (
        deck.map((e, i) => {
          return (
            <Card
              onClick={() => setSelected(e)}
              card={e}
              cl={card(i, overlap)}
              key={i}
              onError={onError}
            />
          )
        })
      ) : (
        <div className={classes} />
      )}
    </div>
  )
}

export default ImageBox
